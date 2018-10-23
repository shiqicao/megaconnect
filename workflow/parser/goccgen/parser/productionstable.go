// Code generated by gocc; DO NOT EDIT.

package parser

import (
    pa "github.com/megaspacelab/megaconnect/workflow/parser/utils"
    wf "github.com/megaspacelab/megaconnect/workflow"    
)

type (
	//TODO: change type and variable names to be consistent with other tables
	ProdTab      [numProductions]ProdTabEntry
	ProdTabEntry struct {
		String     string
		Id         string
		NTType     int
		Index      int
		NumSymbols int
		ReduceFunc func([]Attrib) (Attrib, error)
	}
	Attrib interface {
	}
)

var productionsTable = ProdTab{
	ProdTabEntry{
		String: `S' : Monitor	<<  >>`,
		Id:         "S'",
		NTType:     0,
		Index:      0,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Monitor : kdMonitor id kdChain id kdCondition Expr kdVar "{" VarDecls "}"	<< pa.MonitorAction(X[1], X[3], X[5], X[8]) >>`,
		Id:         "Monitor",
		NTType:     1,
		Index:      1,
		NumSymbols: 10,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.MonitorAction(X[1], X[3], X[5], X[8])
		},
	},
	ProdTabEntry{
		String: `VarDecl : id "=" Expr	<< pa.VarDeclAction(X[0], X[2]) >>`,
		Id:         "VarDecl",
		NTType:     2,
		Index:      2,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.VarDeclAction(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `VarDecls : VarDecl	<<  >>`,
		Id:         "VarDecls",
		NTType:     3,
		Index:      3,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `VarDecls : VarDecls VarDecl	<< pa.VarDeclsAction(X[0], X[1]) >>`,
		Id:         "VarDecls",
		NTType:     3,
		Index:      4,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.VarDeclsAction(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `MonitorProperty : kdChain ":" id	<<  >>`,
		Id:         "MonitorProperty",
		NTType:     4,
		Index:      5,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `MonitorProperty : kdCondition ":" Expr	<<  >>`,
		Id:         "MonitorProperty",
		NTType:     4,
		Index:      6,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `MonitorProperties : MonitorProperty	<<  >>`,
		Id:         "MonitorProperties",
		NTType:     5,
		Index:      7,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `MonitorProperties : MonitorProperty MonitorProperties	<<  >>`,
		Id:         "MonitorProperties",
		NTType:     5,
		Index:      8,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Expr : Expr "==" Term	<< wf.NewBinOp(wf.EqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Expr",
		NTType:     6,
		Index:      9,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.EqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Expr : Expr "!=" Term	<< wf.NewBinOp(wf.NotEqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Expr",
		NTType:     6,
		Index:      10,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.NotEqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Expr : Term	<<  >>`,
		Id:         "Expr",
		NTType:     6,
		Index:      11,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Term : boolLit	<< pa.TermBoolLitAction(X[0]) >>`,
		Id:         "Term",
		NTType:     7,
		Index:      12,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.TermBoolLitAction(X[0])
		},
	},
	ProdTabEntry{
		String: `Term : "(" Expr ")"	<<  >>`,
		Id:         "Term",
		NTType:     7,
		Index:      13,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
}
