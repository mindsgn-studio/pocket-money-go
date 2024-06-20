import { useEffect } from "react"
import { CircularProgress } from "@chakra-ui/react"
import { WalletExists, GetWallets } from "../../wailsjs/go/main/App";
import { useWallet } from "../context";

function Loading() {
    const {isReady, auth, walletExist} = useWallet()

    useEffect(() => {
        if(isReady){
            walletExist()
        }
    },[isReady])

    return (
       <div id="app-loading">
            <CircularProgress isIndeterminate color='black' />
       </div>
    )
}

export default Loading
