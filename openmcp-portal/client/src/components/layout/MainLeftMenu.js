import React, { Component } from "react";
import Tree from "@naisutech/react-tree";

class MainLeftMenu extends Component {
  render() {
    //   ㅇ 메뉴
    // - 모니터링 > 대시보드, 다차원 메트릭스
    // - 멀티플클러스터스 > 클러스터스, 노드
    // - workload > 프로젝트, 디플로이, 파드, 네트워크
    // - 어드벤스드 > 스냅샷, 마이그레이션
    // - 셋팅스 > 설정, 미터링

    const onSelectLeftMenu = (props) => {
      console.log(props);

      if (props.length > 0 && props[0].hasOwnProperty("url")) {
        this.props.info.history.push(props[0].url);
      }
      // this.props.propsData.info.history.push("/nodes");
    };

    const nodes = [
      {
        label: "Monitor",
        id: 1000,
        parentId: null,
        items: [
          {
            label: "Dashboard",
            parentId: 1000,
            id: {
              itemId: 1001,
              url: "/dashboard",
            },
          },
          {
            label: "Multiple Metrics",
            parentId: 1000,
            id: {
              itemId: 1002,
              url: "/nodes",
            },
          },
        ],
      },
      {
        label: "Multiple Clusters",
        id: 2000,
        parentId: null,
        items: [
          {
            label: "Clusters",
            parentId: 2000,
            id: {
              itemId: 2001,
              url: "/clusters",
            },
          },
          {
            label: "Nodes",
            parentId: 2000,
            id: {
              itemId: 2002,
              url: "/nodes",
            },
          },
        ],
      },
      {
        label: "Workloads",
        id: 3000,
        parentId: null,
        items: [
          {
            label: "Projects",
            parentId: 3000,
            id: {
              itemId: 3001,
              url: "/projects",
            },
          },
          {
            label: "Deployments",
            parentId: 3000,
            id: {
              itemId: 3002,
              url: "/deployments",
            },
          },
        ],
      },
      {
        label: "Advenced",
        id: 4000,
        parentId: null,
        items: [
          {
            label: "Snapshots",
            parentId: 4000,
            id: {
              itemId: 4001,
              url: "/maintenance/snapshot",
            },
          },
          {
            label: "Migrations",
            parentId: 4000,
            id: {
              itemId: 4002,
              url: "/maintenance/migration",
            },
            
          },
        ],
      },
      {
        label: "Settings",
        id: 5000,
        parentId: null,
        items: [
          {
            label: "Accounts",
            parentId: 5000,
            id: {
              itemId: 5001,
              url: "/settings/accounts",
            },
          },
          {
            label: "Group Role",
            parentId: 5000,
            id: {
              itemId: 5002,
              url: "/settings/group-role",
            },
          },
          {
            label: "Policy",
            parentId: 5000,
            id: {
              itemId: 5003,
              url: "/settings/accounts",
            },
          },
          {
            label: "Alert",
            parentId: 5000,
            id: {
              itemId: 5004,
              url: "/settings/alert",
            },
          },
          {
            label: "Config",
            parentId: 5000,
            id: {
              itemId: 5005,
              url: "/settings/config",
            },
          },
          
          {
            label: "Meterings",
            parentId: 5000,
            id: {
              itemId: 5006,
              url: "/settings/config",
            },
          },
        ],
      },
    ];

    return (
      <aside className="main-sidebar">
        <section className="sidebar">
          {/* <Tree nodes={nodes}  /> */}
          <Tree nodes={nodes} onSelect={onSelectLeftMenu} />
        </section>
      </aside>
    );
  }
}

export default MainLeftMenu;
