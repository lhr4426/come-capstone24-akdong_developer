package main

import (
    "fmt"
    "time"
)

// 컨테이너를 나타내는 구조체(각 컨테이너 정보 저장)
type Container struct {
    ID     string
    Status string
}

// 노드를 나타내는 구조체(각 노드 정보 저장)
type Node struct {
    ID          string
    Containers  []*Container // 해당 노드에 배치된 컨테이너 목록
    MaxCapacity int // 최대 수용 가능한 컨테이너 개수
}

// 오케스트레이터(컨테이너를 스케줄링하고 노드 관리)
type Orchestrator struct {
    Nodes []*Node // 노드의 목록
}

// 새로운 노드 추가
func (o *Orchestrator) AddNode(node *Node) {
    o.Nodes = append(o.Nodes, node)
}

// 컨테이너 스케줄링
func (o *Orchestrator) ScheduleContainer(container *Container) {
    for _, node := range o.Nodes {
        if len(node.Containers) < node.MaxCapacity {
            container.Status = "Running"
            node.Containers = append(node.Containers, container)
            fmt.Printf("Container %s scheduled on Node %s\n", container.ID, node.ID)
            return
        }
    }
    fmt.Printf("Failed to schedule Container %s: No available nodes\n", container.ID)
}

func main() {
    // 오케스트레이터 생성
    orchestrator := &Orchestrator{}

    // 노드 생성
    node1 := &Node{ID: "Node1", MaxCapacity: 2}
    orchestrator.AddNode(node1)
	fmt.Println("node : ",node1)

    // 컨테이너 생성 및 스케줄링
    container1 := &Container{ID: "Container1", Status: "Pending"}
    orchestrator.ScheduleContainer(container1)
	fmt.Println("container1 : ",container1)

    // 잠시 대기 후 상태 출력
    time.Sleep(1 * time.Second)
    fmt.Printf("Container status: %s\n", container1.Status)
}
