package test

import (
	wf "github.com/megaspacelab/megaconnect/workflow"
)

type NodeVisitor struct {
	VisitNode  func(node wf.Node) interface{}
	VisitExpr  func(node wf.Expr) interface{}
	VisitBinOp func(node wf.BinOp) interface{}
}

func VisitNode(node Node, v NodeVisitor) interface{} {
	v.VisitNode(node)
	switch n := node.(type) {
	case Expr:
		v.VisitExpr(n)
		switch n1 := n.(type) {
			case *wf.BinOp {
				v.VisitBinOp(n1)
				VisiteNode(node1.Left(), v) 
				VisiteNode(node1.Right(), v) 				
			}
		}
	}
}
