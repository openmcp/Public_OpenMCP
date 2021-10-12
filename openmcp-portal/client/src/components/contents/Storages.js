import React, { Component } from "react";
import Paper from "@material-ui/core/Paper";
import { NavLink, Link } from "react-router-dom";
import CircularProgress from "@material-ui/core/CircularProgress";
import {
  SearchState,
  IntegratedFiltering,
  PagingState,
  IntegratedPaging,
  SortingState,
  IntegratedSorting,
} from "@devexpress/dx-react-grid";
import {
  Grid,
  Table,
  Toolbar,
  SearchPanel,
  TableHeaderRow,
  PagingPanel,
} from "@devexpress/dx-react-grid-material-ui";
import Editor from "./../modules/Editor";

class Storages extends Component {
  constructor(props) {
    super(props);
    this.state = {
      columns: [
        { name: "project_name", title: "Name" },
        { name: "project_status", title: "Status" },
        { name: "project_creator", title: "Createor" },
        { name: "project_created_time", title: "Created Time" },
      ],
      
      rows: "",

      // Paging Settings
      currentPage: 0,
      setCurrentPage: 0,
      pageSize: 5,
      setPageSize: 5,
      pageSizes: [5, 10, 15, 0],

      completed: 0,
      editorContext : ``,
    };
  }

  componentWillMount() {
    this.props.menuData("none");
  }

  

  callApi = async () => {
    const response = await fetch("/api/projects");
    const body = await response.json();
    return body;
  };

  progress = () => {
    const { completed } = this.state;
    this.setState({ completed: completed >= 100 ? 0 : completed + 1 });
  };

  //컴포넌트가 모두 마운트가 되었을때 실행된다.
  componentDidMount() {
    //데이터가 들어오기 전까지 프로그래스바를 보여준다.
    this.timer = setInterval(this.progress, 20);
    this.callApi()
      .then((res) => {
        if(res == null){
          this.setState({ rows: [] });
        } else {
          this.setState({ rows: res });
        }
        clearInterval(this.timer);
      })
      .catch((err) => console.log(err));
  }

  render() {
    // 셀 데이터 스타일 변경
    const HighlightedCell = ({ value, style, row, ...restProps }) => (
      <Table.Cell>
        <span
          style={{
            color:
              value === "Healthy"
                ? "#1ab726"
                : value === "Unhealthy"
                ? "red"
                : undefined,
          }}
        >
          {value}
        </span>
      </Table.Cell>
    );

    //셀
    const Cell = (props) => {
      const { column } = props;
      if (column.name === "project_status") {
        return <HighlightedCell {...props} />;
      } else if (column.name === "project_name") {
        return (
          <Table.Cell
            component={Link}
            to={{
              pathname: `/projects/${props.value}/overview`,
              state: {
                test : "testvalue"
              }
            }}
            {...props}
            style={{ cursor: "pointer" }}
          ></Table.Cell>
        );
      }
      return <Table.Cell {...props} />;
    };

    const HeaderRow = ({ row, ...restProps }) => (
      <Table.Row
        {...restProps}
        style={{
          cursor: "pointer",
          backgroundColor: "whitesmoke",
          // ...styles[row.sector.toLowerCase()],
        }}
        // onClick={()=> alert(JSON.stringify(row))}
      />
    );
    // const i = 0;
    const Row = (props) => {
      return <Table.Row {...props} />;
    };

    return (
      <div className="content-wrapper fulled">
        {/* 컨텐츠 헤더 */}
        <section className="content-header">
          <h1>
          Storages
            <small></small>
          </h1>
          <ol className="breadcrumb">
            <li>
              <NavLink to="/dashboard">Home</NavLink>
            </li>
            <li className="active">Storages</li>
          </ol>
        </section>
        <section className="content" style={{ position: "relative" }}>
          <Paper>
            {this.state.rows ? (
              [
                // <input type="button" value="create"></input>,
                <Editor title="create" context={this.state.editorContext}/>,
                <Grid
                  rows={this.state.rows}
                  columns={this.state.columns}
                  style={{ color: "red" }}
                >
                  <Toolbar />
                  {/* 검색 */}
                  <SearchState defaultValue="" />
                  <IntegratedFiltering />
                  <SearchPanel style={{ marginLeft: 0 }} />

                  {/* Sorting */}
                  <SortingState
                  // defaultSorting={[{ columnName: 'city', direction: 'desc' }]}
                  />
                  <IntegratedSorting />

                  {/* 페이징 */}
                  <PagingState defaultCurrentPage={0} defaultPageSize={5} />
                  <IntegratedPaging />
                  <PagingPanel pageSizes={this.state.pageSizes} />

                  {/* 테이블 */}
                  <Table cellComponent={Cell} rowComponent={Row} />
                  <TableHeaderRow
                    showSortingControls
                    rowComponent={HeaderRow}
                  />
                </Grid>,
              ]
            ) : (
              <CircularProgress
                variant="determinate"
                value={this.state.completed}
                style={{ position: "absolute", left: "50%", marginTop: "20px" }}
              ></CircularProgress>
            )}
          </Paper>
        </section>
      </div>
    );
  }
}

export default Storages;
