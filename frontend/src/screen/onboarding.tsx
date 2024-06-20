import NewWalletForm from "../component/form"
import { Box, Heading } from "@chakra-ui/react"

function Onbaording() {
    return (
        <div id="app-home">
            <Box
                display={"flex"}
                flexDir={"column"}
                background="white"
                padding={20}
                paddingX={200}
                borderRadius={10}>
                <Box >
                    <Box>
                        <Heading
                            fontFamily={"SF Rounded Bold"}
                            size={"md"}>POCKET MONEY</Heading>
                    </Box>
                </Box>
                <NewWalletForm />
            </Box>
        </div>
    )
}

export default Onbaording
