import { Flex } from "@chakra-ui/react";
import React, { Children } from "react";
import Navbar from "./Navbar";

const Layout: React.FC = ({ children }) => {
  return (
    <Flex
      direction="column"
      align="center"
      // maxW={{ xl: "1200px" }}
      minHeight="100vh"
      margin="0 auto"
      background="gray.100"
      paddingBottom={16}
    >
      <Navbar />
      {children}
    </Flex>
  );
};

export default Layout;
