package bps_test

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"go.mercari.io/go-bps/bps"
)

func TestOneAmountEquality(t *testing.T) {
	t.Parallel()

	oneAmt := bps.NewFromAmount(1)

	per := bps.NewFromPercentage(100)
	if !oneAmt.Equal(per) {
		t.Error("1 amount = 100%")
	}

	bp := bps.NewFromBasisPoint(10000)
	if !oneAmt.Equal(bp) {
		t.Error("1 amount = 10,000 basis points")
	}

	dbp := bps.NewFromDeciBasisPoint(100000)
	if !oneAmt.Equal(dbp) {
		t.Error("1 amount = 100,000 deci basis points")
	}

	ppm := bps.NewFromPPM(big.NewInt(1000000))
	if !oneAmt.Equal(ppm) {
		t.Error("1 amount = 1000,000 ppm")
	}

	ppb := bps.NewFromPPB(big.NewInt(1000000000))
	if !oneAmt.Equal(ppb) {
		t.Error("1 amount = 1000,000,000 ppb")
	}
}

func TestNewFromString(t *testing.T) {
	tests := map[string]struct {
		arg     string
		want    *bps.BPS
		wantErr bool
	}{
		"int part and decimal part": {
			"123.456",
			bps.NewFromBasisPoint(1234560),
			false,
		},
		"only int part": {
			"123",
			bps.NewFromBasisPoint(1230000),
			false,
		},
		"only decimal part": {
			".1234",
			bps.NewFromBasisPoint(1234),
			false,
		},
		"negative value": {
			"-123.456",
			bps.NewFromBasisPoint(-1234560),
			false,
		},
		"zero": {
			"0.0",
			bps.NewFromAmount(0),
			false,
		},
		"short zero": {
			".0",
			bps.NewFromAmount(0),
			false,
		},
		"If include multi dots, it should return an error": {
			"123.45.6",
			nil,
			true,
		},
		"If base 2 format, it should return an error": {
			"0b11",
			nil,
			true,
		},
		"If base 8 format, it should return an error": {
			"0o75",
			nil,
			true,
		},
		"If base 16 format, it should return an error": {
			"0xF5",
			nil,
			true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := bps.NewFromString(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustFromString(t *testing.T) {
	tests := map[string]struct {
		arg       string
		want      *bps.BPS
		wantPanic bool
	}{
		"int part and decimal part": {
			"123.456",
			bps.NewFromBasisPoint(1234560),
			false,
		},
		"only int part": {
			"123",
			bps.NewFromBasisPoint(1230000),
			false,
		},
		"only decimal part": {
			".1234",
			bps.NewFromBasisPoint(1234),
			false,
		},
		"negative value": {
			"-123.456",
			bps.NewFromBasisPoint(-1234560),
			false,
		},
		"zero": {
			"0.0",
			bps.NewFromAmount(0),
			false,
		},
		"short zero": {
			".0",
			bps.NewFromAmount(0),
			false,
		},
		"If include multi dots, it should return an error": {
			"123.45.6",
			nil,
			true,
		},
		"If base 2 format, it should return an error": {
			"0b11",
			nil,
			true,
		},
		"If base 8 format, it should return an error": {
			"0o75",
			nil,
			true,
		},
		"If base 16 format, it should return an error": {
			"0xF5",
			nil,
			true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tt.wantPanic {
				//nolint:gocritic
				defer func() {
					err := recover()
					if err == nil {
						t.Error("MustFromString() should occure a panic")
					}
				}()
			}
			got := bps.MustFromString(tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleNewFromString() {
	// 15%
	b1, _ := bps.NewFromString("0.15")
	fmt.Println(b1.Percentages())
	// 2.645% = 264.5 basis points
	b2, _ := bps.NewFromString("0.02645")
	fmt.Println(b2.DeciBasisPoints())
	// Output:
	// 15
	// 2645
}

func ExampleMustFromString() {
	// 15%
	b1 := bps.MustFromString("0.15")
	fmt.Println(b1.Percentages())
	// 2.645% = 264.5 basis points
	b2 := bps.MustFromString("0.02645")
	fmt.Println(b2.DeciBasisPoints())

	a := bps.NewFromAmount(1e12)
	b, _ := bps.NewFromString(".000001") // 1 / 1e6
	// Set PPM as BaseUnit to show value as ppm
	bps.BaseUnit = bps.PPM
	fmt.Println(a.Add(b), "ppm")
	// teardown
	bps.BaseUnit = bps.DeciBasisPoint

	n := bps.NewFromAmount(0)
	for i := 0; i < 1000; i++ {
		n = n.Add(bps.MustFromString(".01"))
	}
	fmt.Println(n.Amounts())
	// Output:
	// 15
	// 2645
	// 1000000000000000001 ppm
	// 10
}

func ExampleNewFromBaseUnit() {
	// backup
	u := bps.BaseUnit
	var arg int64 = 15

	// The default BaseUnit is DeciBasisPoint
	deci := bps.NewFromBaseUnit(arg)
	fmt.Println(deci.PPMs())

	// BaseUnit is updated by PPB
	bps.BaseUnit = bps.PPB
	ppb := bps.NewFromBaseUnit(arg)
	fmt.Println(ppb.PPBs())

	// BaseUnit is updated by PPM
	bps.BaseUnit = bps.PPM
	ppm := bps.NewFromBaseUnit(arg)
	fmt.Println(ppm.PPBs())

	// BaseUnit is updated by HalfBasisPoint
	bps.BaseUnit = bps.HalfBasisPoint
	hbp := bps.NewFromBaseUnit(arg)
	fmt.Println(hbp.PPBs())

	// BaseUnit is updated by BasisPoint
	bps.BaseUnit = bps.BasisPoint
	bp := bps.NewFromBaseUnit(arg)
	fmt.Println(bp.PPBs())

	// BaseUnit is updated by Percentage
	bps.BaseUnit = bps.Percentage
	p := bps.NewFromBaseUnit(arg)
	fmt.Println(p.PPBs())

	// teardown
	bps.BaseUnit = u
	// Output:
	// 150
	// 15
	// 15000
	// 750000
	// 1500000
	// 150000000
}
