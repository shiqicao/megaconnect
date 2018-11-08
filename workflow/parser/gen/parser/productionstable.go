// Code generated by gocc; DO NOT EDIT.

package parser

import (
    pa "github.com/megaspacelab/megaconnect/workflow/parser/utils"
    wf "github.com/megaspacelab/megaconnect/workflow"    
    "github.com/megaspacelab/megaconnect/workflow/parser/gen/token"
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
		String: `S' : Workflow	<<  >>`,
		Id:         "S'",
		NTType:     0,
		Index:      0,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Workflow : kdWorkflow id "{" Decls "}"	<< wf.NewWorkflowDecl(pa.Id(X[1]), 0).AddChildren(X[3].([]wf.Decl)), nil >>`,
		Id:         "Workflow",
		NTType:     1,
		Index:      1,
		NumSymbols: 5,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewWorkflowDecl(pa.Id(X[1]), 0).AddChildren(X[3].([]wf.Decl)), nil
		},
	},
	ProdTabEntry{
		String: `Decls : Decl	<< []wf.Decl{X[0].(wf.Decl)}, nil >>`,
		Id:         "Decls",
		NTType:     2,
		Index:      2,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return []wf.Decl{X[0].(wf.Decl)}, nil
		},
	},
	ProdTabEntry{
		String: `Decls : Decls Decl	<< append(X[0].([]wf.Decl), X[1].([]wf.Decl)[0]), nil >>`,
		Id:         "Decls",
		NTType:     2,
		Index:      3,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return append(X[0].([]wf.Decl), X[1].([]wf.Decl)[0]), nil
		},
	},
	ProdTabEntry{
		String: `Decl : Monitor	<<  >>`,
		Id:         "Decl",
		NTType:     3,
		Index:      4,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Decl : Event	<<  >>`,
		Id:         "Decl",
		NTType:     3,
		Index:      5,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Decl : empty	<<  >>`,
		Id:         "Decl",
		NTType:     3,
		Index:      6,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `Monitor : kdMonitor id kdChain id kdCondition Expr MonitorVar kdFire id ObjLit	<< pa.MonitorAction(X[1], X[3], X[5], X[6], X[8], X[9]) >>`,
		Id:         "Monitor",
		NTType:     4,
		Index:      7,
		NumSymbols: 10,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.MonitorAction(X[1], X[3], X[5], X[6], X[8], X[9])
		},
	},
	ProdTabEntry{
		String: `MonitorVar : kdVar "{" VarDecls "}"	<< X[2], nil >>`,
		Id:         "MonitorVar",
		NTType:     5,
		Index:      8,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[2], nil
		},
	},
	ProdTabEntry{
		String: `MonitorVar : empty	<<  >>`,
		Id:         "MonitorVar",
		NTType:     5,
		Index:      9,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `Event : kdEvent id "{" ObjFields "}"	<< wf.NewEventDecl(pa.Lit(X[1]), wf.NewObjType(X[3].(wf.IdToTy))), nil >>`,
		Id:         "Event",
		NTType:     6,
		Index:      10,
		NumSymbols: 5,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewEventDecl(pa.Lit(X[1]), wf.NewObjType(X[3].(wf.IdToTy))), nil
		},
	},
	ProdTabEntry{
		String: `ObjField : id ":" Type	<< []interface{}{X[0], X[2]}, nil >>`,
		Id:         "ObjField",
		NTType:     7,
		Index:      11,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return []interface{}{X[0], X[2]}, nil
		},
	},
	ProdTabEntry{
		String: `ObjFields_ : ObjFields_ "," ObjField	<< pa.ObjFieldsAction(X[0], X[2]) >>`,
		Id:         "ObjFields_",
		NTType:     8,
		Index:      12,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.ObjFieldsAction(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `ObjFields_ : empty	<< wf.NewIdToTy(), nil >>`,
		Id:         "ObjFields_",
		NTType:     8,
		Index:      13,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewIdToTy(), nil
		},
	},
	ProdTabEntry{
		String: `ObjFields : ObjField ObjFields_	<< pa.ObjFieldsAction(X[1], X[0]) >>`,
		Id:         "ObjFields",
		NTType:     9,
		Index:      14,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.ObjFieldsAction(X[1], X[0])
		},
	},
	ProdTabEntry{
		String: `ObjFields : empty	<< wf.NewIdToTy(), nil >>`,
		Id:         "ObjFields",
		NTType:     9,
		Index:      15,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewIdToTy(), nil
		},
	},
	ProdTabEntry{
		String: `Type : kdStr	<< wf.StrType, nil >>`,
		Id:         "Type",
		NTType:     10,
		Index:      16,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.StrType, nil
		},
	},
	ProdTabEntry{
		String: `Type : kdInt	<< wf.IntType, nil >>`,
		Id:         "Type",
		NTType:     10,
		Index:      17,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.IntType, nil
		},
	},
	ProdTabEntry{
		String: `Type : kdBool	<< wf.BoolType, nil >>`,
		Id:         "Type",
		NTType:     10,
		Index:      18,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.BoolType, nil
		},
	},
	ProdTabEntry{
		String: `Type : "{" ObjFields "}"	<< wf.NewObjType(X[1].(wf.IdToTy)), nil >>`,
		Id:         "Type",
		NTType:     10,
		Index:      19,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewObjType(X[1].(wf.IdToTy)), nil
		},
	},
	ProdTabEntry{
		String: `ObjLit : "{" ObjLitFields "}"	<< X[1], nil >>`,
		Id:         "ObjLit",
		NTType:     11,
		Index:      20,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[1], nil
		},
	},
	ProdTabEntry{
		String: `ObjLitField : id ":" Expr	<< []interface{}{X[0], X[2]}, nil >>`,
		Id:         "ObjLitField",
		NTType:     12,
		Index:      21,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return []interface{}{X[0], X[2]}, nil
		},
	},
	ProdTabEntry{
		String: `ObjLitFields_ : ObjLitFields_ "," ObjLitField	<< pa.ObjLitFieldsAction(X[0], X[2]) >>`,
		Id:         "ObjLitFields_",
		NTType:     13,
		Index:      22,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.ObjLitFieldsAction(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `ObjLitFields_ : empty	<< wf.NewIdToExpr(), nil >>`,
		Id:         "ObjLitFields_",
		NTType:     13,
		Index:      23,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewIdToExpr(), nil
		},
	},
	ProdTabEntry{
		String: `ObjLitFields : ObjLitField ObjLitFields_	<< pa.ObjLitFieldsAction(X[1], X[0]) >>`,
		Id:         "ObjLitFields",
		NTType:     14,
		Index:      24,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.ObjLitFieldsAction(X[1], X[0])
		},
	},
	ProdTabEntry{
		String: `ObjLitFields : empty	<< wf.NewIdToExpr(), nil >>`,
		Id:         "ObjLitFields",
		NTType:     14,
		Index:      25,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewIdToExpr(), nil
		},
	},
	ProdTabEntry{
		String: `VarDecl : id "=" Expr	<< pa.VarDeclAction(X[0], X[2]) >>`,
		Id:         "VarDecl",
		NTType:     15,
		Index:      26,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.VarDeclAction(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `VarDecl : empty	<<  >>`,
		Id:         "VarDecl",
		NTType:     15,
		Index:      27,
		NumSymbols: 0,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return nil, nil
		},
	},
	ProdTabEntry{
		String: `VarDecls : VarDecl	<<  >>`,
		Id:         "VarDecls",
		NTType:     16,
		Index:      28,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `VarDecls : VarDecls VarDecl	<< pa.VarDeclsAction(X[0], X[1]) >>`,
		Id:         "VarDecls",
		NTType:     16,
		Index:      29,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.VarDeclsAction(X[0], X[1])
		},
	},
	ProdTabEntry{
		String: `Expr : Expr "||" Term1	<< wf.NewBinOp(wf.OrOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Expr",
		NTType:     17,
		Index:      30,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.OrOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Expr : Expr "&&" Term1	<< wf.NewBinOp(wf.AndOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Expr",
		NTType:     17,
		Index:      31,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.AndOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Expr : Term1	<<  >>`,
		Id:         "Expr",
		NTType:     17,
		Index:      32,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Term1 : Term1 "==" Term2	<< wf.NewBinOp(wf.EqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term1",
		NTType:     18,
		Index:      33,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.EqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term1 : Term1 "!=" Term2	<< wf.NewBinOp(wf.NotEqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term1",
		NTType:     18,
		Index:      34,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.NotEqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term1 : Term1 ">" Term2	<< wf.NewBinOp(wf.GreaterThanOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term1",
		NTType:     18,
		Index:      35,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.GreaterThanOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term1 : Term1 ">=" Term2	<< wf.NewBinOp(wf.GreaterThanEqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term1",
		NTType:     18,
		Index:      36,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.GreaterThanEqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term1 : Term1 "<" Term2	<< wf.NewBinOp(wf.LessThanOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term1",
		NTType:     18,
		Index:      37,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.LessThanOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term1 : Term1 "<=" Term2	<< wf.NewBinOp(wf.LessThanEqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term1",
		NTType:     18,
		Index:      38,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.LessThanEqualOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term1 : Term2	<<  >>`,
		Id:         "Term1",
		NTType:     18,
		Index:      39,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Term2 : Term2 "+" Term3	<< wf.NewBinOp(wf.PlusOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term2",
		NTType:     19,
		Index:      40,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.PlusOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term2 : Term2 "-" Term3	<< wf.NewBinOp(wf.MinusOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term2",
		NTType:     19,
		Index:      41,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.MinusOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term2 : Term3	<<  >>`,
		Id:         "Term2",
		NTType:     19,
		Index:      42,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Term3 : Term3 "*" Term4	<< wf.NewBinOp(wf.MultOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term3",
		NTType:     20,
		Index:      43,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.MultOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term3 : Term3 "/" Term4	<< wf.NewBinOp(wf.DivOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil >>`,
		Id:         "Term3",
		NTType:     20,
		Index:      44,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewBinOp(wf.DivOp, X[0].(wf.Expr), X[2].(wf.Expr)), nil
		},
	},
	ProdTabEntry{
		String: `Term3 : Term4	<<  >>`,
		Id:         "Term3",
		NTType:     20,
		Index:      45,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Term4 : Term4 "." id	<< pa.ObjAccessorAction(X[0], X[2]) >>`,
		Id:         "Term4",
		NTType:     21,
		Index:      46,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.ObjAccessorAction(X[0], X[2])
		},
	},
	ProdTabEntry{
		String: `Term4 : Term5	<<  >>`,
		Id:         "Term4",
		NTType:     21,
		Index:      47,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	ProdTabEntry{
		String: `Term5 : intLit	<< pa.IntLitAction(X[0]) >>`,
		Id:         "Term5",
		NTType:     22,
		Index:      48,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.IntLitAction(X[0])
		},
	},
	ProdTabEntry{
		String: `Term5 : boolLit	<< pa.BoolLitAction(X[0]) >>`,
		Id:         "Term5",
		NTType:     22,
		Index:      49,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.BoolLitAction(X[0])
		},
	},
	ProdTabEntry{
		String: `Term5 : stringLit	<< pa.StrLitAction(X[0]) >>`,
		Id:         "Term5",
		NTType:     22,
		Index:      50,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.StrLitAction(X[0])
		},
	},
	ProdTabEntry{
		String: `Term5 : ObjLit	<< pa.ObjLitAction(X[0]) >>`,
		Id:         "Term5",
		NTType:     22,
		Index:      51,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return pa.ObjLitAction(X[0])
		},
	},
	ProdTabEntry{
		String: `Term5 : id	<< wf.NewVar(string(X[0].(*token.Token).Lit)), nil >>`,
		Id:         "Term5",
		NTType:     22,
		Index:      52,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return wf.NewVar(string(X[0].(*token.Token).Lit)), nil
		},
	},
	ProdTabEntry{
		String: `Term5 : "(" Expr ")"	<< X[1], nil >>`,
		Id:         "Term5",
		NTType:     22,
		Index:      53,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[1], nil
		},
	},
}
