import React, { Component } from "react";
import { withStyles } from "@material-ui/core/styles";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import MuiDialogTitle from "@material-ui/core/DialogTitle";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import Typography from "@material-ui/core/Typography";
import DialogActions from "@material-ui/core/DialogActions";
import DialogContent from "@material-ui/core/DialogContent";
import Slider from "@material-ui/core/Slider";
import * as utilLog from "../../../util/UtLogs.js";
import { AsyncStorage } from "AsyncStorage";
import axios from "axios";

import Grid from "@material-ui/core/Grid";
import List from "@material-ui/core/List";
import Card from "@material-ui/core/Card";
import CardHeader from "@material-ui/core/CardHeader";
import ListItem from "@material-ui/core/ListItem";
import ListItemText from "@material-ui/core/ListItemText";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import Checkbox from "@material-ui/core/Checkbox";
import Divider from "@material-ui/core/Divider";

function not(a, b) {
  return a.filter((value) => b.indexOf(value) === -1);
}

function intersection(a, b) {
  return a.filter((value) => b.indexOf(value) !== -1);
}

function union(a, b) {
  return [...a, ...not(b, a)];
}

const styles = (theme) => ({
  root: {
    margin: 0,
    padding: theme.spacing(2),
  },
  closeButton: {
    position: "absolute",
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
});

class DashboardSelectModule extends Component {
  constructor(props) {
    super(props);
    this.state = {
      checked : [],
      left : ["system_status","cluster_groups","openmcp_management_clusters","world_cluster_status"],
      right : []
    }
  }

  componentWillMount() {}

  onChange(e) {}

  handleClickOpen = () => {
    this.setState({
      open: true
    });
  };

  handleClose = () => {
    this.setState({ open: false });
  };

  
  handleSave = (e) => {
    //DB에 대시보드 표시 항목을 저장
    const url = `/settings/policy/openmcp-policy`;
    const data = {
      policyId: this.state.policyId,
      rate: {
        start: this.state.g_rate_value[0],
        end: this.state.g_rate_value[1],
      },
      period: {
        start: this.state.period_value[0],
        end: this.state.period_value[1],
      },
    };

    axios
      .put(url, data)
      .then((res) => {
        console.log("res", res.data);
        if (res.data.data.rowCount > 0) {
          // log - policy update
        } else {
          this.props.onUpdateData();
          // console.log("sdfsdf",this.props)
          let userId = null;
          AsyncStorage.getItem("userName", (err, result) => {
            userId = result;
          });
          utilLog.fn_insertPLogs(userId, "log-PO-MD01");
          alert(res.data.message);
        }
      })
      .catch((err) => {
        alert(err);
      });

    let userId = null;
    AsyncStorage.getItem("userName", (err, result) => {
      userId = result;
    });
    utilLog.fn_insertPLogs(userId, "log-PD-MD01");
    this.setState({ open: false });
  };

  render() {
    const DialogTitle = withStyles(styles)((props) => {
      const { children, classes, onClose, ...other } = props;
      return (
        <div>
          <MuiDialogTitle disableTypography className={classes.root} {...other}>
            <Typography variant="h6">{children}</Typography>
            {onClose ? (
              <IconButton
                aria-label="close"
                className={classes.closeButton}
                onClick={onClose}
              >
                <CloseIcon />
              </IconButton>
            ) : null}
          </MuiDialogTitle>
        </div>
      );
    });

    const leftChecked = intersection(this.state.checked, this.state.left);
    const rightChecked = intersection(this.state.checked, this.state.right);

    const handleToggle = (value) => () => {
      const currentIndex = this.state.checked.indexOf(value);
      const newChecked = [...this.state.checked];

      if (currentIndex === -1) {
        newChecked.push(value);
      } else {
        newChecked.splice(currentIndex, 1);
      }

      this.setState({checked : newChecked});
    };

    const numberOfChecked = (items) => intersection(this.state.checked, items).length;

    const handleToggleAll = (items) => () => {
      if (numberOfChecked(items) === items.length) {
        this.setState({checked : not(this.state.checked, items)});
      } else {
        this.setState({checked : union(this.state.checked, items)});
      }
    };

    const handleCheckedRight = () => {
      this.setState({right : this.state.right.concat(leftChecked)});
      this.setState({left : not(this.state.left, leftChecked)});
      this.setState({checked : not(this.state.checked, leftChecked)});
    };

    const handleCheckedLeft = () => {
      this.setState({left : this.state.left.concat(rightChecked) });
      this.setState({right : not(this.state.right, rightChecked)});
      this.setState({checked : not(this.state.checked, rightChecked)});
    };

    const customList = (title, items) => (
      <Card className="ddd" style={{minWidth:"350px", minHeight:"350px"}}>
        <CardHeader
          sx={{ px: 2, py: 1 }}
          avatar={
            <Checkbox
              onClick={handleToggleAll(items)}
              checked={
                numberOfChecked(items) === items.length && items.length !== 0
              }
              indeterminate={
                numberOfChecked(items) !== items.length &&
                numberOfChecked(items) !== 0
              }
              disabled={items.length === 0}
              inputProps={{
                "aria-label": "all items selected",
              }}
            />
          }
          title={title}
          subheader={`${numberOfChecked(items)}/${items.length} selected`}
        />
        <Divider />
        <List
          sx={{
            width: 200,
            height: 230,
            bgcolor: "background.paper",
            overflow: "auto",
          }}
          dense
          component="div"
          role="list"
        >
          {items.map((value) => {
            const labelId = `transfer-list-all-item-${value}-label`;

            return (
              <ListItem
                key={value}
                role="listitem"
                button
                onClick={handleToggle(value)}
              >
                <ListItemIcon>
                  <Checkbox
                    checked={this.state.checked.indexOf(value) !== -1}
                    tabIndex={-1}
                    disableRipple
                    inputProps={{
                      "aria-labelledby": labelId,
                    }}
                  />
                </ListItemIcon>
                <ListItemText id={labelId} primary={`${value}`} />
              </ListItem>
            );
          })}
          <ListItem />
        </List>
      </Card>
    );

    return (
      <div style={{display:"inline-block"}}>
        <Button
          variant="outlined"
          color="primary"
          onClick={this.handleClickOpen}
          style={{
            position: "absolute",
            left: "150px",
            top: "45px",
            // zIndex: "10",
            width: "148px",
            height: "31px",
            textTransform: "capitalize",
          }}
        >
          Edit Dashboard
        </Button>
        <Dialog
          onClose={this.handleClose}
          aria-labelledby="customized-dialog-title"
          open={this.state.open}
          fullWidth={false}
          maxWidth={false}
        >
          <DialogTitle id="customized-dialog-title" onClose={this.handleClose}>
            Dashboard Custom View
          </DialogTitle>
          <DialogContent dividers>
            <Grid
              container
              spacing={2}
              justifyContent="center"
              alignItems="center"
            >
              <Grid item>{customList("Choices", this.state.left)}</Grid>
              <Grid item>
                <Grid container direction="column" alignItems="center">
                  <Button
                    sx={{ my: 0.5 }}
                    variant="outlined"
                    size="small"
                    onClick={handleCheckedRight}
                    disabled={leftChecked.length === 0}
                    aria-label="move selected right"
                  >
                    &gt;
                  </Button>
                  <Button
                    sx={{ my: 0.5 }}
                    variant="outlined"
                    size="small"
                    onClick={handleCheckedLeft}
                    disabled={rightChecked.length === 0}
                    aria-label="move selected left"
                  >
                    &lt;
                  </Button>
                </Grid>
              </Grid>
              <Grid item>{customList("Chosen", this.state.right)}</Grid>
            </Grid>
          </DialogContent>
          <DialogActions>
            <Button onClick={this.handleSave} color="primary">
              save
            </Button>
            <Button onClick={this.handleClose} color="primary">
              cancel
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

export default DashboardSelectModule;
