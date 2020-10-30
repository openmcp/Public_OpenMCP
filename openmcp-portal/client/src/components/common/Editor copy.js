// import React, { Component } from "react";
// import AceEditor from "react-ace";
// import "ace-builds/src-noconflict/mode-yaml";
// import "ace-builds/src-noconflict/theme-nord_dark";

// import { withStyles } from "@material-ui/core/styles";
// import Button from "@material-ui/core/Button";
// import Dialog from "@material-ui/core/Dialog";
// import MuiDialogTitle from "@material-ui/core/DialogTitle";
// import MuiDialogContent from "@material-ui/core/DialogContent";
// import MuiDialogActions from "@material-ui/core/DialogActions";
// import IconButton from "@material-ui/core/IconButton";
// import CloseIcon from "@material-ui/icons/Close";
// import Typography from "@material-ui/core/Typography";

// // import Modal from "@material-ui/core/Modal";

// const styles = (theme) => ({
//   root: {
//     margin: 0,
//     padding: theme.spacing(2),
//   },
//   closeButton: {
//     position: "absolute",
//     right: theme.spacing(1),
//     top: theme.spacing(1),
//     color: theme.palette.grey[500],
//   },
// });

// class Editor extends Component {
//   constructor(props) {
//     super(props);
//     console.log("state")
//     this.state = {
//       open: false,
//     };
//   }

//   render() {
//     const DialogTitle = withStyles(styles)((props) => {
//       const { children, classes, onClose, ...other } = props;
//       return (
//         <MuiDialogTitle disableTypography className={classes.root} {...other}>
//           <Typography variant="h6">{children}</Typography>
//           {onClose ? (
//             <IconButton
//               aria-label="close"
//               className={classes.closeButton}
//               onClick={onClose}
//             >
//               <CloseIcon />
//             </IconButton>
//           ) : null}
//         </MuiDialogTitle>
//       );
//     });

//     const DialogContent = withStyles((theme) => ({
//       root: {
//         padding: theme.spacing(2),
//       },
//     }))(MuiDialogContent);

//     const DialogActions = withStyles((theme) => ({
//       root: {
//         margin: 0,
//         padding: theme.spacing(1),
//       },
//     }))(MuiDialogActions);

//     const onChange = (newValue) => {
//       console.log("change", newValue);
//     };

//     const handleClickOpen = () => {
//       this.setState({ open: true });
//     };

//     const handleClose = () => {
//       this.setState({ open: false });
//     };

//     // const useStyles = makeStyles((theme) => ({
//     //     paper: {
//     //       position: 'absolute',
//     //       width: 400,
//     //       backgroundColor: theme.palette.background.paper,
//     //       border: '2px solid #000',
//     //       boxShadow: theme.shadows[5],
//     //       padding: theme.spacing(2, 4, 3),
//     //     },
//     //   }));

//     const body2 = (
//       <div>
//         <Button variant="outlined" color="primary" onClick={handleClickOpen}>
//           Open dialog
//         </Button>
//         <Dialog
//           onClose={handleClose}
//           aria-labelledby="customized-dialog-title"
//           open={this.state.open}
//         >
//           <DialogTitle id="customized-dialog-title" onClose={handleClose}>
//             Modal title
//           </DialogTitle>
//           <DialogContent dividers>
//             {/* <Typography gutterBottom>
//             Cras mattis consectetur purus sit amet fermentum. Cras justo odio, dapibus ac facilisis
//             in, egestas eget quam. Morbi leo risus, porta ac consectetur ac, vestibulum at eros.
//           </Typography>
//           <Typography gutterBottom>
//             Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Vivamus sagittis
//             lacus vel augue laoreet rutrum faucibus dolor auctor.
//           </Typography>
//           <Typography gutterBottom>
//             Aenean lacinia bibendum nulla sed consectetur. Praesent commodo cursus magna, vel
//             scelerisque nisl consectetur et. Donec sed odio dui. Donec ullamcorper nulla non metus
//             auctor fringilla.
//           </Typography> */}
//           </DialogContent>
//           <DialogActions>
//             <Button autoFocus onClick={handleClose} color="primary">
//               Save changes
//             </Button>
//           </DialogActions>
//         </Dialog>
//       </div>
//     );

//     const editor = (
//       <AceEditor
//         mode="yaml"
//         theme="nord_dark"
//         onChange={onChange}
//         name="UNIQUE_ID_OF_DIV"
//         editorProps={{ $blockScrolling: true }}
//         value={`apiVersion: types.kubefed.io/v1beta1
//           kind: FederatedDeployment
//           metadata:
//             annotations:
//               kubesphere.io/creator: admin
//             labels:
//               app: mapp
//             name: mapp
//             namespace: mpr
//           spec:
//             overrides:
//               - clusterName: host
//                 clusterOverrides:
//                   - path: /spec/replicas
//                     value: 3
//               - clusterName: slave
//                 clusterOverrides:
//                   - path: /spec/replicas
//                     value: 3
//             placement:
//               clusters:
//                 - name: host
//                 - name: slave
//             template:
//               metadata:
//                 annotations:
//                   kubesphere.io/containerSecrets: null
//                 labels:
//                   app: mapp
//                 namespace: mpr
//               spec:
//                 replicas: 1
//                 selector:
//                   matchLabels:
//                     app: mapp
//                 strategy:
//                   rollingUpdate:
//                     maxSurge: 25%
//                     maxUnavailable: 25%
//                   type: RollingUpdate
//                 template:
//                   metadata:
//                     labels:
//                       app: mapp
//                   spec:
//                     affinity: {}
//                     containers:
//                       - image: 'paulbouwer/hello-kubernetes:1.8'
//                         imagePullPolicy: IfNotPresent
//                         name: container-aft8xf
//                         ports:
//                           - containerPort: 8080
//                             name: http-8080
//                             protocol: TCP
//                     imagePullSecrets: null
//                     initContainers: []
//                     serviceAccount: default
//                     volumes: []
//                     `}
//       />
//     );
//     console.log("render")

//     return (
//       <div>
//         {/* <button type="button" onClick={handleOpen}>
//           Open Modal
//         </button>

//         <Modal
//           open={this.state.open}
//           onClose={handleClose}
//           aria-labelledby="simple-modal-title"
//           aria-describedby="simple-modal-description"
//         >
//           {body}
//         </Modal> */}
//         <Button variant="outlined" color="primary" onClick={handleClickOpen}>
//           Open dialog
//         </Button>
//         <Dialog
//           onClose={handleClose}
//           aria-labelledby="customized-dialog-title"
//           open={this.state.open}
//         >
//           <DialogTitle id="customized-dialog-title" onClose={handleClose}>
//             Modal title
//           </DialogTitle>
//           <DialogContent dividers>
//             {/* <Typography gutterBottom>
//             Cras mattis consectetur purus sit amet fermentum. Cras justo odio, dapibus ac facilisis
//             in, egestas eget quam. Morbi leo risus, porta ac consectetur ac, vestibulum at eros.
//           </Typography>
//           <Typography gutterBottom>
//             Praesent commodo cursus magna, vel scelerisque nisl consectetur et. Vivamus sagittis
//             lacus vel augue laoreet rutrum faucibus dolor auctor.
//           </Typography>
//           <Typography gutterBottom>
//             Aenean lacinia bibendum nulla sed consectetur. Praesent commodo cursus magna, vel
//             scelerisque nisl consectetur et. Donec sed odio dui. Donec ullamcorper nulla non metus
//             auctor fringilla.
//           </Typography> */}
//           {editor}
//           </DialogContent>
//           <DialogActions>
//             <Button autoFocus onClick={handleClose} color="primary">
//               Save changes
//             </Button>
//           </DialogActions>
//         </Dialog>
//       </div>
//     );
//   }
// }

// export default withStyles(styles)(Editor);
