export function getMenu(pathParams) {
  const menuData = {
    clusters : [
      {
        type: "single",
        title : "OverView",
        path : "/clusters/" + pathParams + "/overview",
      },
      {
        type: "single",
        title : "Nodes",
        path : "/clusters/" + pathParams + "/nodes",
      },
      {
        type: "single",
        title : "Pods",
        path : "/clusters/" + pathParams + "/pods",
      },
      {
        type: "single",
        title : "Storage Class",
        path : "/clusters/" + pathParams + "/storage_class",
      },
      {
        type: "multi",
        title : "Settings",
        path : "/clusters/" + pathParams + "/settings",
        sub : [
          { title: "Workloads", path: "/clusters/"+pathParams+"/settings/member" },
        ]
      }
    ],
    nodes : [
      {
        type: "single",
        title : "OverView",
        path : "/nodes/" + pathParams + "/overview",
      },
      {
        type: "multi",
        title : "Resources",
        path : "/nodes/" + pathParams + "/resources",
        sub : [
          { title: "Workloads", path: "/nodes/"+pathParams+"/resources/workloads" },
          { title: "Pods", path: "/nodes/"+pathParams+"/resources/pods" },
          { title: "Services", path: "/nodes/"+pathParams+"/resources/services" },
          { title: "Ingress", path: "/nodes/"+pathParams+"/resources/ingress" },
        ]
      }
    ],
    projects : [
      {
        type: "single",
        title : "OverView",
        path : "/projects/" + pathParams + "/overview",
      },
      {
        type: "multi",
        title : "Resources",
        path : "/projects/" + pathParams + "/resources",
        sub : [
          { title: "Workloads", path: "/projects/"+pathParams+"/resources/workloads" },
          { title: "Pods", path: "/projects/"+pathParams+"/resources/pods" },
          { title: "Services", path: "/projects/"+pathParams+"/resources/services" },
          { title: "Ingress", path: "/projects/"+pathParams+"/resources/ingress" },
        ]
      },
      {
        type: "single",
        title : "Volumes",
        path : "/projects/" + pathParams + "/volumes",
      },
      {
        type: "multi",
        title : "Config",
        path : "/projects/" + pathParams + "/config",
        sub : [
          { title: "Secrets", path: "/projects/" + pathParams + "/config/secrets"},
          { title: "ConfigMaps", path: "/projects/"+pathParams+"/config/configMaps" },
        ]
      },
      {
        type: "multi",
        title : "Settings",
        path : "/projects/" + pathParams + "/settings",
        sub : [
          { title: "Members", path: "/projects/" + pathParams + "/settings/members"},
        ]
      }
    ],
    pods : [
      {
        type: "single",
        title : "OverView",
        path : "/pods/" + pathParams + "/overview",
      },
      {
        type: "multi",
        title : "Resources",
        path : "/nodes/" + pathParams + "/resources",
        sub : [
          { title: "Workloads", path: "/pods/"+pathParams+"/resources/workloads" },
          { title: "Pods", path: "/pods/"+pathParams+"/resources/pods" },
          { title: "Services", path: "/pods/"+pathParams+"/resources/services" },
          { title: "Ingress", path: "/pods/"+pathParams+"/resources/ingress" },
        ]
      }
    ]
  }
  return menuData;
}