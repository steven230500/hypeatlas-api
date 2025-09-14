package entities

import "gorm.io/gorm"

type Relations []Relation

type Relation struct {
	Name       string
	Conditions []interface{}
	Joins      []string
	Where      []string
	Functions  []func(db *gorm.DB) *gorm.DB
}

func (r Relations) Preload(prefix string, db *gorm.DB) *gorm.DB {
	for _, rel := range r {
		db = rel.Preload(prefix, db)
		db = rel.Join(db)
		db = rel.FilterWhere(db)
		db = rel.ApplyFunctions(db)
	}
	return db
}

func (r Relation) Preload(prefix string, db *gorm.DB) *gorm.DB {
	return db.Preload(prefix+r.Name, r.Conditions...)
}

func (r Relation) Join(db *gorm.DB) *gorm.DB {
	for _, j := range r.Joins {
		db = db.Joins(j)
	}
	return db
}

func (r Relation) FilterWhere(db *gorm.DB) *gorm.DB {
	for _, w := range r.Where {
		db = db.Where(w)
	}
	return db
}

func (r Relation) ApplyFunctions(db *gorm.DB) *gorm.DB {
	for _, f := range r.Functions {
		db = f(db)
	}
	return db
}
