package model

import (
	. "app_api/database/mysql"
	"app_api/util/log"
	"time"
)

type InvoiceList struct {
	Id            int64     `table:"id"`
	InvoiceNumber string    `table:"invoice_number"`
	MachineNo     string    `table:"machine_no"`
	Year          int       `table:"year"`
	Month         int       `table:"month"`
	InvoiceType   int       `table:"type"`
	Status        int       `table:"status"`
	UpdatedAt     time.Time `table:"updated_at"`
	CreatedAt     time.Time `table:"created_at"`
}

func (model *InvoiceList) SetInvoiceNumber(invoiceNumber string) *InvoiceList {
	model.InvoiceNumber = invoiceNumber
	return model
}

func (model *InvoiceList) SetMachineNo(machineNo string) *InvoiceList {
	model.MachineNo = machineNo
	return model
}

func (model *InvoiceList) SetYear(year int) *InvoiceList {
	model.Year = year
	return model
}

func (model *InvoiceList) SetMonth(month int) *InvoiceList {
	model.Month = month
	return model
}

func (model *InvoiceList) SetInvoiceType(invoiceType int) *InvoiceList {
	model.InvoiceType = invoiceType
	return model
}

func (model *InvoiceList) SetStatus(status int) *InvoiceList {
	model.Status = status
	return model
}

func (model *InvoiceList) Create() (int64, error) {
	model.CreatedAt = time.Now()
	return Model(model).Insert()
}

func (model *InvoiceList) Update(columns []string) error {
	return Model(model).Where("id", "=", model.Id).Update(columns)
}

func (model *InvoiceList) QueryOne() *InvoiceList {
	table := Model(model)

	if model.Status > 0 {
		table.Where("status", "=", model.Status)
	}

	if model.MachineNo != "" {
		table.Where("machine_no", "=", model.MachineNo)
	}

	log.Error(table.Select([]string{"id", "invoice_number"}).
		OrderBy([]string{"id"}, []string{"asc"}).
		Find().Scan(&model.Id, &model.InvoiceNumber))
	return model
}

func (model *InvoiceList) Count() int {
	table := Model(model)

	if model.Status > 0 {
		table.Where("status", "=", model.Status)
	}

	if model.MachineNo != "" {
		table.Where("machine_no", "=", model.MachineNo)
	}

	return table.Count()
}
