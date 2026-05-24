package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	mysql_db "saidakbar.origin/db/mysql"
	mysql_entity "saidakbar.origin/db/mysql/entity"
	"saidakbar.origin/repository/models"
)

type TemplateMysqlRepository interface {
	CreateTemplete(templete *mysql_entity.TempleteModel) (err error)
	UpdateTemplete(templete *mysql_entity.TempleteModel) (err error)
	DeleteTemplete(templete *mysql_entity.TempleteModel) (err error)
	GetTemplete(id string, organizationID string) (templete *mysql_entity.TempleteModel, err error)
	GetAllTempletes(filters *models.GetAllTempleteFilter) (templetes []*mysql_entity.TempleteModel, err error)
	GetTempleteCount(filters *models.GetTempleteCountFilter) (count int64, err error)
}

type templatMysqlRepository struct {
	db *gorm.DB
}

func NewTemplateMysqlRepository(db mysql_db.Database) TemplateMysqlRepository {
	return &templatMysqlRepository{db: db.GetDB()}
}

func (r *templatMysqlRepository) CreateTemplete(templete *mysql_entity.TempleteModel) (err error) {
	return r.db.Create(templete).Error
}

func (r *templatMysqlRepository) UpdateTemplete(templete *mysql_entity.TempleteModel) (err error) {
	return r.db.Model(templete).Updates(map[string]interface{}{
		mysql_entity.FieldTempleteName:      templete.Name,
		mysql_entity.FieldTempleteUpdatedAt: templete.UpdatedAt,
	}).Error
}

func (r *templatMysqlRepository) DeleteTemplete(templete *mysql_entity.TempleteModel) (err error) {
	return r.db.Model(templete).Updates(map[string]interface{}{
		mysql_entity.FieldTempleteIsDeleted: true,
		mysql_entity.FieldTempleteDeletedAt: templete.DeletedAt,
	}).Error
}

func (r *templatMysqlRepository) GetTemplete(id string, organizationID string) (templete *mysql_entity.TempleteModel, err error) {
	var entity mysql_entity.TempleteModel
	err = r.db.Scopes(
		entity.FilterID(id),
		entity.FilterOrganizationID(organizationID),
		entity.FilterIsDeleted(false),
	).First(templete).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *templatMysqlRepository) GetAllTempletes(filters *models.GetAllTempleteFilter) (templetes []*mysql_entity.TempleteModel, err error) {
	var entity mysql_entity.TempleteModel

	scopes := []func(*gorm.DB) *gorm.DB{
		entity.FilterOrganizationID(filters.OrganizationID),
		entity.FilterIsDeleted(false),
		entity.OrderByCreatedAtDesc(),
	}

	if filters.Search != nil && strings.TrimSpace(*filters.Search) != "" {
		scopes = append(scopes, entity.FilterName(*filters.Search))
	}
	if filters.StartDate != nil {
		scopes = append(scopes, entity.FilterCreatedAtGte(*filters.StartDate))
	}
	if filters.EndDate != nil {
		scopes = append(scopes, entity.FilterCreatedAtLte(*filters.EndDate))
	}

	query := r.db.Model(&entity).Scopes(scopes...)
	if filters.Skip != nil {
		query = query.Offset(*filters.Skip)
	}
	if filters.Limit != nil {
		query = query.Limit(*filters.Limit)
	}

	err = query.Find(&templetes).Error
	return
}

func (r *templatMysqlRepository) GetTempleteCount(filters *models.GetTempleteCountFilter) (count int64, err error) {
	var entity mysql_entity.TempleteModel

	scopes := []func(*gorm.DB) *gorm.DB{
		entity.FilterOrganizationID(filters.OrganizationID),
		entity.FilterIsDeleted(false),
	}

	if filters.Search != nil && strings.TrimSpace(*filters.Search) != "" {
		scopes = append(scopes, entity.FilterName(*filters.Search))
	}
	if filters.StartDate != nil {
		scopes = append(scopes, entity.FilterCreatedAtGte(*filters.StartDate))
	}
	if filters.EndDate != nil {
		scopes = append(scopes, entity.FilterCreatedAtLte(*filters.EndDate))
	}

	err = r.db.Model(&entity).Scopes(scopes...).Count(&count).Error
	return
}
