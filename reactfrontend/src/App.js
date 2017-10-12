import 'rc-calendar/assets/index.css';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import Calendar from 'rc-calendar';
import logo from './logo.svg';
import './App.css';

import { Button } from 'semantic-ui-react'

//export default ButtonExampleButton


class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Felis Resource Management System</h1>
        </header>
            <Button>Login</Button>
        {/*<p className="App-intro">
          To get started, edit <code>src/App.js</code> and save to reload.
         </p>*/}

         <div class="ui celled grid">
         <div class="row">
           <div class="three wide column">
               <p>Find available resource today.</p>
               <Calendar />
               
           </div>
           <div class="thirteen wide column">
             <img class="ui image" src="/assets/images/wireframe/centered-paragraph.png" />
           </div>
         </div>         
       </div>
        
      </div>
    );
  }
}

export default App;
