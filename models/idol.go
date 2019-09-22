package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Idol struct {
	Ability      string `orm:"column(ability);size(255);null"`
	Age          string `orm:"column(age);size(255);null"`
	Association  string `orm:"column(association);size(255);null"`
	Birthday     string `orm:"column(birthday);size(255);null"`
	Bloodtype    string `orm:"column(bloodtype);size(255);null"`
	Cv           string `orm:"column(cv);size(255);null"`
	DominantHand string `orm:"column(dominant_hand);size(255);null"`
	Favorite     string `orm:"column(favorite);size(255);null"`
	Height       string `orm:"column(height);size(255);null"`
	Hiragana     string `orm:"column(hiragana);size(255);null"`
	Hobby        string `orm:"column(hobby);size(255);null"`
	Hometown     string `orm:"column(hometown);size(255);null"`
	Id           int    `orm:"column(id);pk"`
	Mlid         int    `orm:"column(mlid);null"`
	Name         string `orm:"column(name);size(255);null"`
	Size         string `orm:"column(size);size(255);null"`
	Type         string `orm:"column(type);size(255);null"`
	Weight       string `orm:"column(weight);size(255);null"`
}

func (t *Idol) TableName() string {
	return "idol"
}

func init() {
	orm.RegisterModel(new(Idol))
}

// AddIdol insert a new Idol into database and returns
// last inserted Id on success.
func AddIdol(m *Idol) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetIdolById retrieves Idol by Id. Returns error if
// Id doesn't exist
func GetIdolById(id int) (v *Idol, err error) {
	o := orm.NewOrm()
	v = &Idol{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllIdol retrieves all Idol matches certain condition. Returns empty list if
// no records exist
func GetAllIdol(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Idol))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Idol
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateIdol updates Idol by Id and returns error if
// the record to be updated doesn't exist
func UpdateIdolById(m *Idol) (err error) {
	o := orm.NewOrm()
	v := Idol{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteIdol deletes Idol by Id and returns error if
// the record to be deleted doesn't exist
func DeleteIdol(id int) (err error) {
	o := orm.NewOrm()
	v := Idol{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Idol{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
