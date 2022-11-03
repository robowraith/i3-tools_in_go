package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Node struct {
	ID int64           `json:"id"`
	Parent *Node
	Orientation string `json:"orientation"`
	Type string        `json:"type"`
	Focused bool       `json:"focused"`
	Children []Node    `json:"Nodes"`
}

func focus_left() {
    exec.Command("i3-msg", "focus", "left").Run()
}

func focus_right() {
    exec.Command("i3-msg", "focus", "right").Run()
}

func workspace_next() {
    exec.Command("i3-msg", "workspace", "next").Run()
}

func get_layout() Node {
	layout := Node{}

	layout_json, err := exec.Command("i3-msg", "-t", "get_tree").Output()
	if err != nil {
        fmt.Println(err)
    }

	_ = json.Unmarshal([]byte(layout_json), &layout)

	return layout
}

func find_focused_window(node Node, parent Node) Node {
	node.Parent = &parent
	if node.Focused == false {
		for _, child := range (node.Children) {
			current_node := find_focused_window(child, node)
			if current_node.Focused == false {
				continue
			} else {
				return current_node
			}
		}
	}
	return node
}

func is_left_most_window(window Node) bool {
	parent := *window.Parent
    if window.Type == "workspace" {
		return true
	} else if parent.Orientation == "horizontal" && window.ID != parent.Children[0].ID {
		return false
	}
	return is_left_most_window(parent)
}

func is_right_most_window(window Node) bool {
	parent := *window.Parent
    if window.Type == "workspace" {
		return true
	} else if parent.Orientation == "horizontal" && window.ID != parent.Children[len(parent.Children)-1].ID {
		return false
	}
	return is_right_most_window(parent)
}

func go_all_the_way_left() {
	parent := Node {
	    Orientation: "",
	}

	focused_window := find_focused_window(get_layout(), parent)

	if ! is_left_most_window(focused_window) {
		focus_left()
		go_all_the_way_left()
	}
}

func main()  {
	parent := Node {
	    Orientation: "",
	}

	focused_window := find_focused_window(get_layout(), parent)

	if is_right_most_window(focused_window) == true {
		workspace_next()
		go_all_the_way_left()
	} else {
		focus_right()
	}
}
