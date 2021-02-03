import React from "react";
import { BrowserRouter, Switch, Route } from "react-router-dom";

import Layout from "./components/Layout";
import CreateCheck from "./pages/CreateCheck";
import Dashboard from "./pages/Dashboard";
import CheckDetails from "./pages/CheckDetails";

const App = () => {
  return (
    <BrowserRouter>
      <Layout>
        <Switch>
          <Route exact path="/dashboard" component={Dashboard} />
          <Route exact path="/new" component={CreateCheck} />
          <Route path="/check/:id" component={CheckDetails} />
        </Switch>
      </Layout>
    </BrowserRouter>
  );
};

export default App;
