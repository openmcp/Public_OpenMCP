// UtLogs is group of functions what writing portal users behaviors
// import React, { Component } from "react";
import axios from 'axios';

export function fn_insertPLogs(userid, code) {
  const url = `/apimcp/portal-log`;
      const data = {
        userid:userid,
        code:code,
      };
  axios.post(url, data)
      .then((res) => {
        // console.log()
      })
      .catch((err) => {
          // console.log()
      });
}