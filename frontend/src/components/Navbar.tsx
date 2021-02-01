import React, { useState } from "react";
import { Text, Flex, Heading, Box, Link } from "@chakra-ui/react";
import { Link as RouterLink } from "react-router-dom";
import ThemeToggler from "../ThemeToggler";
import { HamburgerIcon } from "@chakra-ui/icons";

const MenuItems: React.FC = ({ children }) => (
  <Text mt={{ base: 4, md: 0 }} mr={6} display="block">
    {children}
  </Text>
);

const Navbar: React.FC = (props) => {
  const [show, setShow] = useState(false);
  const handleToggle = () => setShow(!show);

  return (
    <Flex
      as="nav"
      align="center"
      justify="space-between"
      wrap="wrap"
      p={4}
      mb={16}
      w="100%"
      // bg={["white", "white", "transparent", "transparent"]}
      // color={["white", "white", "gray.900", "gray.900"]}
      // bg="white"
      color="gray.900"
      // boxShadow="xl"
      {...props}
    >
      <Flex align="center">
        <Link as={RouterLink} to="/dashboard">
          <Heading as="h1" size="lg" color="teal.500" letterSpacing={"-.1rem"}>
            WebMonitor
          </Heading>
        </Link>
      </Flex>

      {/* <Spacer /> */}
      {/* <Box display={{ base: "block", md: "none" }} onClick={handleToggle}>
        <HamburgerIcon xlinkTitle="Menu" />
      </Box> */}

      {/* <Box
        display={{ sm: show ? "block" : "none", md: "block" }}
        width={{ sm: "full", md: "auto" }}
        alignItems="center"
      > */}
      <Flex align="center">
        <MenuItems>
          <ThemeToggler />
        </MenuItems>
      </Flex>
      {/* </Box> */}
    </Flex>
  );
};

export default Navbar;
