package cmd

import "context"

type Parameter struct {
	ParamsId int64
	VtnId    int64
	VtpId    int64
	Times    int64
	Signals  string
}

func (p Parameter) InsertQuery() string {
	return "insert or ignore into Parameter(VtnId, VtpId, Times, Signals) values (?,?,?,?)"
}

func (p Parameter) Insert(ctx context.Context, repository Repository) error {
	db, err := repository.Connect()
	defer db.Db.Close()
	if err != nil {
		return err
	}

	_, err = db.WithContext(ctx).Exec(p.InsertQuery(), p.VtnId, p.VtpId, p.Times, p.Signals)
	return err
}

func (p Parameter) SelectQuery() string {
	return "select * from Parameter where VtnId = ? and VtpId = ? and Times = ? and Signals = ?"
}

func (p Parameter) Select(ctx context.Context, repository Repository) (Parameter, error) {
	db, err := repository.Connect()
	defer db.Db.Close()
	if err != nil {
		return Parameter{}, err
	}

	var rt Parameter
	err = db.WithContext(ctx).SelectOne(&rt, p.SelectQuery(), p.VtnId, p.VtpId, p.Times, p.Signals)
	return rt, err
}

func (t Task) NewParameter(ctx context.Context, r Repository) (Parameter, error) {
	var rt Parameter

	var VtnId, VtpId int64

	// insert and select id from Transistor Table for Vtn
	{
		err := t.Vtn.Insert(ctx, r)
		if err != nil {
			return rt, err
		}
		rec, err := t.Vtn.Select(ctx, r)
		VtnId = rec.TransistorId
		if err != nil {
			return rt, err
		}
	}

	// insert and select id from Transistor Table for Vtp
	{
		err := t.Vtp.Insert(ctx, r)
		rec, err := t.Vtp.Select(ctx, r)
		VtpId = rec.TransistorId
		if err != nil {
			return rt, err
		}
	}

	// get json-string for Signals
	json, err := t.PlotPoint.ToJson()
	if err != nil {
		return rt, err
	}

	// insert and select id from Parameter Table for param

	p := Parameter{
		VtnId:   VtnId,
		VtpId:   VtpId,
		Times:   int64(t.Times),
		Signals: json,
	}

	err = p.Insert(ctx, r)
	if err != nil {
		return rt, err
	}

	rec, err := p.Select(ctx, r)
	if err != nil {
		return rt, err
	}
	return rec, nil
}
