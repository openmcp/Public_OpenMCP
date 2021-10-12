import React, { Component } from "react";
import MenuItem from "@material-ui/core/MenuItem";
import TranslateIcon from "@material-ui/icons/Translate";
import IconButton from "@material-ui/core/IconButton";
import Popper from "@material-ui/core/Popper";
import MenuList from "@material-ui/core/MenuList";
import Grow from "@material-ui/core/Grow";
import Paper from "@material-ui/core/Paper";
import ClickAwayListener from '@material-ui/core/ClickAwayListener';

// import List from "@material-ui/core/List";
// import ListItem from "@material-ui/core/ListItem";
// import ListItemText from "@material-ui/core/ListItemText";
// import Menu from "@material-ui/core/Menu";


class LanguageListMenu extends Component {
  constructor(props) {
    super(props);
    this.state = {
      anchorEl: null,
      selectedIndex: "1",
    };
  }
  render() {
    const open = Boolean(this.state.anchorEl);
    const handleClickListItem = (event) => {
      this.setState({ anchorEl: event.currentTarget });
    };

    const handleMenuItemClick = (event, index) => {
      this.setState({ selectedIndex: index });
      this.setState({ anchorEl: null });
    };

    const handleClose = () => {
      this.setState({ anchorEl: null });
    };

    const options = ["English", "한국어"];

    return (
      <div className="language">
        <IconButton
          aria-label="more"
          aria-controls="long-menu"
          aria-haspopup="true"
          onClick={handleClickListItem}
          className="lang-list-btn"
        >
          <TranslateIcon className="language-icon" />
        </IconButton>
        <Popper
          open={open}
          anchorEl={this.state.anchorEl}
          role={undefined}
          transition
          disablePortal
          placement={"bottom-end"}
        >
          {({ TransitionProps, placement }) => (
            <Grow
              {...TransitionProps}
              style={{
                transformOrigin:
                  placement === "bottom" ? "center top" : "center top",
              }}
            >
              <Paper>
                <ClickAwayListener onClickAway={handleClose}>
                  <MenuList autoFocusItem={open} id="menu-list-grow">
                    {options.map((option, index) => (
                      <MenuItem
                        key={option}
                        // disabled={index === 0}
                        selected={index === this.state.selectedIndex}
                        onClick={(event) => handleMenuItemClick(event, index)}
                        className="lang-list-menu-item"
                      >
                        {option}
                      </MenuItem>
                    ))}
                  </MenuList>
                </ClickAwayListener>
              </Paper>
            </Grow>
          )}
        </Popper>
      </div>
    );
  }
}

export default LanguageListMenu;
