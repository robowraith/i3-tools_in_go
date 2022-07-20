package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Node struct {
	Parent *Node
	Orientation string `json:"orientation"`
	Type string        `json:"type"`
	Focused bool       `json:"focused"`
	Nodes []Node       `json:"nodes"`
}

func get_layout() Node {
	fmt.Println("get_layout")
	layout := Node{}
	
	layout_json, err := exec.Command("i3-msg", "-t", "get_tree").Output()
	if err != nil {
        fmt.Println(err)
    }
	
	_ = json.Unmarshal([]byte(layout_json), &layout)
	
	return layout
}

func find_focused_window(node Node, parent Node) Node {
	fmt.Printf("find_focused_window(%v, %v)\n", node.Type, node.Focused)
    node.Parent = &parent
	if node.Focused == false {
		for _, node := range (node.Nodes) {
			node := find_focused_window(node, parent)
			if node.Focused == false {
				continue
			} else {
				return node
			}
		}
	}
	return node
}

func is_left_most_window(window Node) bool {
	fmt.Printf("is_left_most_window(%v %v)\n", window.Type, window.Focused)
	parent := window.Parent
    if window.Type == "workspace" {
		return true
	} else if parent.Orientation == "horizontal" && &window != &parent.Nodes[0] {
		return false
	}
	return is_left_most_window(*parent)
}

func is_right_most_window(window Node) bool {
    if window.Type == "workspace" {
		return true
	}
	parent := window.Parent
	if parent.Orientation == "horizontal" && &window == &parent.Nodes[len(parent.Nodes)-1] {
		return true
	}
	return is_right_most_window(*parent)
}

func go_all_the_way_right() {
	parent := Node {
	    Orientation: "",
	}
    
	focused_window := find_focused_window(get_layout(), parent)
	
	if ! is_right_most_window(focused_window) {
		focus_right()
		go_all_the_way_right()
	}	
}

func focus_left() {
    exec.Command("i3-msg", "-t", "focus", "left")
}

func focus_right() {
    exec.Command("i3-msg", "-t", "focus", "right")
}

func workspace_prev() { 
    exec.Command("i3-msg", "-t", "workspace", "prev")
}

func main()  {
	parent := Node {
	    Orientation: "",
	}

	focused_window := find_focused_window(get_layout(), parent)

	if is_left_most_window(focused_window) {
		workspace_prev()
		go_all_the_way_right()
	} else {
		focus_left()
	}
}
