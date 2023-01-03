package main

import (
  "bytes"
  "fmt"
  "log"
  "go/parser"
  "go/ast"
  "go/token"
  "go/printer"
  "strings"
  "flag"
  "strconv"
  //"encoding/json"
  //"os"


  //"github.com/goccy/go-graphviz"
)

type CodeWalker struct {
    NodesVisited []map[string][]string
    FuncDeclVisited []map[string][]string
    AssignStmtVisited []map[string][]string
    ForStmtVisited []map[string][]string
    IfStmtVisited []map[string][]string
    Fset *token.FileSet
    stdout_mode bool
}

var GLOBAL_WALKER CodeWalker

func (walker CodeWalker) Visit(node_in ast.Node) (visitor_out ast.Visitor) {
	//if node_in == nil { 
	//	visitor_out = nil
	//	return
    //}
    visitor_out = walker
    this_node_map := make(map[string][]string)
    this_node_string := get_node_string(walker.Fset, node_in)
    if len(this_node_string) == 0 {
        return
    }
    this_node_map["node_string"] = this_node_string
    //str_printed := fmt.Sprintf("walker.Dummy2: %v", walker.Dummy2)
    outstring := ""
    node_inFuncDecl, isFuncDecl := node_in.(*ast.FuncDecl)
    if isFuncDecl {
        var return_type_list []string
        //fmt.Printf("node_inFuncDec.Type: %v\n", get_node_string(walker.Fset, node_inFuncDecl.Type))
        if node_inFuncDecl.Type.Results == nil || node_inFuncDecl.Type.Results.List == nil {
            this_node_map["return_types"] = []string{"none"}
            } else {
            for _, return_field := range node_inFuncDecl.Type.Results.List{
                return_type_list = append(return_type_list, get_node_string(walker.Fset, return_field.Type)...)
            }
            this_node_map["return_types"] = return_type_list
        }
        this_node_map["type"] = []string{"FuncDecl"}
        this_node_map["name"] = get_node_string(walker.Fset, node_inFuncDecl.Name)
        this_node_map["body"] = get_node_string(walker.Fset, node_inFuncDecl.Body.List)
        this_node_map["pos"] = []string{fmt.Sprintf("%v",node_inFuncDecl.Pos())}
        this_node_map["end"] = []string{fmt.Sprintf("%v",node_inFuncDecl.End())}
        GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
        GLOBAL_WALKER.FuncDeclVisited = append(GLOBAL_WALKER.FuncDeclVisited, this_node_map)
        if walker.stdout_mode {
          outstring = fmt.Sprintf("node_inFuncDecl.Name: %v\nnode_inFuncDecl.Body.List: %v\n",get_node_string(walker.Fset, node_inFuncDecl.Name),get_node_string(walker.Fset, node_inFuncDecl.Body.List))
          fmt.Println(outstring)
        }
    }
    node_inGenDecl, isGenDecl := node_in.(*ast.GenDecl)
    if isGenDecl {
        this_node_map["type"] = []string{"GenDecl"}
        this_node_map["token"] = get_node_string(walker.Fset, node_inGenDecl.Tok)
        this_node_map["pos"] = []string{fmt.Sprintf("%v",node_inGenDecl.Pos())}//get_node_string(walker.Fset, node_inGenDecl.Pos())}
        this_node_map["end"] = []string{fmt.Sprintf("%v",node_inGenDecl.End())}
        GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
        if walker.stdout_mode {
            fmt.Println("^^^ GENERAL DECLARATION!! vvv")
            outstring = fmt.Sprintf("node_inGenDecl.Tok: %v\n",get_node_string(walker.Fset, node_inGenDecl.Tok))
            fmt.Println(outstring)
        }
        
    }
    node_inIf, isIf := node_in.(*ast.IfStmt)
    if isIf {
        this_node_map["type"] = []string{"IfStmt"}
        this_node_map["condition"] = get_node_string(walker.Fset, node_inIf.Cond)
        this_node_map["body"] = get_node_string(walker.Fset, node_inIf.Body.List)
        this_node_map["pos"] = get_node_string(walker.Fset, node_inIf.If)
        this_node_map["end"] = []string{fmt.Sprintf("%v",node_inIf.End())}
        GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
        GLOBAL_WALKER.IfStmtVisited = append(GLOBAL_WALKER.IfStmtVisited, this_node_map)
        if walker.stdout_mode {
            fmt.Println("^^^ IF STATEMENT!! vvv")
            outstring = fmt.Sprintf("node_inIf.Cond: %v\nnode_inIf.Body.List: %v\n",get_node_string(walker.Fset, node_inIf.Cond),get_node_string(walker.Fset, node_inIf.Body.List))
            fmt.Println(outstring)
        }
        
    }
    node_inGo, isGo := node_in.(*ast.GoStmt)
    if isGo {
        this_node_map["type"] = []string{"GoStmt"}
        this_node_map["call.fun"] = get_node_string(walker.Fset, node_inGo.Call.Fun)
        this_node_map["pos"] = get_node_string(walker.Fset, node_inGo.Go)
        this_node_map["end"] = []string{fmt.Sprintf("%v",node_inGo.End())}
        GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
        if walker.stdout_mode {
            fmt.Println("^^^ GO STATEMENT!! vvv")
            outstring = fmt.Sprintf("node_inGo.Call.Fun: %v\n",get_node_string(walker.Fset, node_inGo.Call.Fun))
            fmt.Println(outstring)
        }
        
    }
    node_inAss, isAss := node_in.(*ast.AssignStmt)
    if isAss {
        this_node_map["type"] = []string{"AssignStmt"}
        this_node_map["pos"] = []string{fmt.Sprintf("%v",node_inAss.Pos())}
        this_node_map["end"] = []string{fmt.Sprintf("%v",node_inAss.End())}
        this_node_map["tok"] = []string{fmt.Sprintf("%v",node_inAss.Tok)}
        var lhs_slice []string
        var rhs_slice []string
        for _, lhs := range(node_inAss.Lhs){
            for _, rhs := range(node_inAss.Rhs){
                //outstring = fmt.Sprintf("%vnode_inAss.Lhs[%v]: %v\nnode_inAss.Rhs[%v]: %v\n",outstring, lhs_ind,get_node_string(walker.Fset, lhs),rhs_ind,get_node_string(walker.Fset, rhs))
                //this_rhs_key := fmt.Sprintf("rhs%v",rhs_ind)
                rhs_slice = append(rhs_slice, get_node_string(walker.Fset, rhs)...)
            }
            //this_lhs_key := fmt.Sprintf("lhs%v",lhs_ind)
            lhs_slice = append(lhs_slice, get_node_string(walker.Fset, lhs)...)
        }
        this_node_map["lhs"] = lhs_slice
        this_node_map["rhs"] = rhs_slice
        GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
        GLOBAL_WALKER.AssignStmtVisited = append(GLOBAL_WALKER.AssignStmtVisited, this_node_map)
        if walker.stdout_mode {
            fmt.Println("^^^ ASSIGNMENT STATEMENT!! vvv")
            outstring = ""
            for lhs_ind, lhs := range(node_inAss.Lhs){
                for rhs_ind, rhs := range(node_inAss.Rhs){
                    outstring = fmt.Sprintf("%vnode_inAss.Lhs[%v]: %v\nnode_inAss.Rhs[%v]: %v\n",outstring, lhs_ind,get_node_string(walker.Fset, lhs),rhs_ind,get_node_string(walker.Fset, rhs))
                }
            }
            fmt.Println(outstring)
        }
    }
    node_inExpr, isExpr := node_in.(*ast.ExprStmt)
    if isExpr {
        this_node_map["type"] = []string{"ExprStmt"}
        this_node_map["expression"] = get_node_string(walker.Fset, node_inExpr.X)
        this_node_map["pos"] = []string{fmt.Sprintf("%v",node_inExpr.Pos())}
        this_node_map["end"] = []string{fmt.Sprintf("%v",node_inExpr.End())}
        GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
        if walker.stdout_mode {
            fmt.Println("^^^ EXPRESSION STATEMENT!! vvv")
            outstring = fmt.Sprintf("node_inExpr.X: %v\n",get_node_string(walker.Fset, node_inExpr.X))
            fmt.Println(outstring)
        }
    }
    node_inFor, isFor := node_in.(*ast.ForStmt)
    if isFor {
        this_node_map["type"] = []string{"ForStmt"}
        this_node_map["condition"] = get_node_string(walker.Fset, node_inFor.Cond)
        this_node_map["body"] = get_node_string(walker.Fset, node_inFor.Body.List)
        this_node_map["pos"] = []string{fmt.Sprintf("%v",node_inFor.Pos())}
        this_node_map["end"] = []string{fmt.Sprintf("%v",node_inFor.End())}
        GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
        GLOBAL_WALKER.ForStmtVisited = append(GLOBAL_WALKER.ForStmtVisited, this_node_map)
        if walker.stdout_mode {
            fmt.Println("^^^ FOR STATEMENT!! vvv")
            outstring = fmt.Sprintf("node_inFor.Cond: %v\nnode_inFor.Body.List: %v\n",get_node_string(walker.Fset, node_inFor.Cond),get_node_string(walker.Fset, node_inFor.Body.List))
            fmt.Println(outstring)
        }
    }
    node_inReturn, isReturn := node_in.(*ast.ReturnStmt)
    if isReturn {
        this_node_map["type"] = []string{"ReturnStmt"}
        this_node_map["pos"] = []string{fmt.Sprintf("%v",node_inReturn.Pos())}
        this_node_map["end"] = []string{fmt.Sprintf("%v",node_inReturn.End())}
        //str_builder := ""
        var results_slice []string
        for _, result := range(node_inReturn.Results){
        //  str_builder = fmt.Sprintf("%vnode_inReturn.Results[%v]: %v\n",str_builder,n,get_node_string(walker.Fset, result))
            //this_key := fmt.Sprintf("results%v",n)
            results_slice = append(results_slice, get_node_string(walker.Fset, result)...)
        }
        this_node_map["results"] = results_slice
        GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
        if walker.stdout_mode {
            fmt.Println("^^^ RETURN STATEMENT!! vvv")
            str_builder := ""
            for n, result := range(node_inReturn.Results){
                str_builder = fmt.Sprintf("%vnode_inReturn.Results[%v]: %v\n",str_builder,n,get_node_string(walker.Fset, result))
            }
            fmt.Println(str_builder)
        }
    }
    node_inSwitch, isSwitch := node_in.(*ast.SwitchStmt)
    if isSwitch {
        this_node_map["type"] = []string{"SwitchStmt"}
        this_node_map["tag"] = get_node_string(walker.Fset, node_inSwitch.Tag)
        this_node_map["body"] = get_node_string(walker.Fset, node_inSwitch.Body.List)
        this_node_map["pos"] = get_node_string(walker.Fset, node_inSwitch.Pos())
        this_node_map["end"] = []string{fmt.Sprintf("%v",node_inSwitch.End())}
        GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
        if walker.stdout_mode {
            fmt.Println("^^^ SWITCH STATEMENT!! vvv")
            outstring = fmt.Sprintf("node_inSwitch.Tag: %v\nnode_inSwitch.Body.List: %v\n",get_node_string(walker.Fset, node_inSwitch.Tag),get_node_string(walker.Fset, node_inSwitch.Body.List))
            fmt.Println(outstring)
        }
    }
    walker.NodesVisited = append(walker.NodesVisited, this_node_map)
    //GLOBAL_WALKER.NodesVisited = append(GLOBAL_WALKER.NodesVisited, this_node_map)
    if walker.stdout_mode {
        _, err := fmt.Printf("walked node\n")
        if err!=nil {
           log.Println(err.Error())
           return
        }
    }
    
    return

}
func main() {
  var filename_in string
  var exports_only bool
  var position_in int
  flag.StringVar(&filename_in, "file", "analyze_me.go", "name of source file to analyze")
  flag.BoolVar(&exports_only, "exports_only", false, "Boolean value - whether to filter all but export nodes or not.")
  flag.IntVar(&position_in, "pos", 1, "Position of the cursor in the input file.")
  flag.Parse()
//  fmt.Println("filename_in: ", filename_in)
//  fmt.Println("filename_out: ", filename_out)
//  fmt.Println("exports_only: ", exports_only)
  fset := token.NewFileSet()
  node, err := parser.ParseFile(fset, filename_in, nil, parser.ParseComments)
  if err != nil {
      log.Fatal(err)
  }
  if exports_only {
      ast.FileExports(node)
  }

  var nodes_list []map[string][]string 
  da_walker := CodeWalker{NodesVisited: nodes_list, Fset: fset, stdout_mode: false}
  ast.Walk(da_walker,node)

//  fmt.Printf("da_walker: %v", da_walker)

  max_range := 0 
  var node_string_list_temp []string
  for _ , func_block := range GLOBAL_WALKER.FuncDeclVisited {
//    fmt.Printf("func_block: %v\n", func_block)
    start, err := strconv.Atoi(func_block["pos"][0])
    if err!=nil {
       log.Fatal(err)
    }
    end, err := strconv.Atoi(func_block["end"][0])
    if err!=nil {
       log.Fatal(err)
    }
//    fmt.Printf("name, start , end : %v , %v , %v\n" , func_block["name"][0], start, end)
    if start < position_in && end > position_in && end - start > max_range {
      node_string_list_temp = func_block["node_string"]
      //fmt.Println(func_block["name"])
      //to find largest satisfactory block in case of nested function declarations:
      max_range = end - start
    }
  }
  str_builder := ""
  for _, line := range node_string_list_temp {
    str_builder = fmt.Sprintf("%v%v\n",str_builder,line)
  }

  fmt.Println(str_builder)

  //fmt.Println(buf.String())
//  marshalled_global_walker, err := json.MarshalIndent(GLOBAL_WALKER,"","    ")
//  if err != nil {
//      log.Fatal(err)
//  }
//  
//  err = os.WriteFile(filename_out, marshalled_global_walker, 0644)
//  if err != nil {
//      log.Fatal(err)
//  }
//  if da_walker.stdout_mode {
//    fmt.Println("marshalled_global_walker: ", string(marshalled_global_walker))
//  }
}

func get_node_string(fset *token.FileSet, node interface{}) (string_list_out []string){
	var thingbuf bytes.Buffer
	printer.Fprint(&thingbuf, fset, node)
	string_out := thingbuf.String()
    string_list_out = strings.Split(string_out, "\n")
	//string_list_out = strings.TrimSpace(strings.ReplaceAll(string_list_out, "\n\t", "\n"))
	return
}

