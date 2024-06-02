import './style.css';
import './App.css';
import { ChakraProvider } from '@chakra-ui/react'
import {
    createBrowserRouter,
    RouterProvider,
} from "react-router-dom";

import Loading from './screen/loading';
import Home from './screen/home';
import Verify from './screen/verify';
import Onboarding from "./screen/onboarding"
import { WalletProvider } from './context';

const router = createBrowserRouter([
    {
        path: "/",
        element: <Loading />,
    },
    {
        path: "/home",
        element: <Home />,
    },
    {
        path: "/onboarding",
        element: <Onboarding />,
    },
    {
        path: "/verify",
        element: <Verify />,
    },
]);

function App() {
    return (
        <WalletProvider>
            <ChakraProvider>
                <RouterProvider router={router} />
            </ChakraProvider>
        </WalletProvider>
    )
}

export default App
