import React, {
    createContext,
    ReactNode,
    useContext,
    useEffect,
  } from 'react';
import { WalletExists } from "../../wailsjs/go/main/App"
import { useCookies } from 'react-cookie'
  
  interface Wallet {
    isReady: boolean;
    auth: boolean;
    walletExist: () => void,
    createNewWallet: () => void,
    createCode: () => void,
    verifyCode: () => void,
  }
  
  const WalletContext = createContext<Wallet>({
    isReady: false,
    auth: false,
    walletExist: () => false,
    createNewWallet: () => {},
    createCode: () => {},
    verifyCode: () => {}
  });
  
  function useWallet(): any {
    const context = useContext(WalletContext);
    if (!context) {
      throw new Error('usePlayer must be used within an PlayerProvider');
    }
    return context;
  }
  
  const WalletProvider = (props: {children: ReactNode}): any => {  
    const [cookies, setCookie] = useCookies(['user'])
    const [isReady, setIsReady] = React.useState(false)
    const [auth, setAuth] = React.useState(false)
    const getCookie = () => {
      
    }

    const walletExist = () => {
        try{
          WalletExists()
          .then((bool: boolean)=>{
            setIsReady(bool)
            console.log(cookies)
            if(bool){
              // window.location.href = "/verify"
            }else{
              // window.location.href = "/onboarding"
            }
          })
          .catch((error: any)=>{
            // window.location.href = "/error"
          })
        } catch(error) {
          // window.location.href = "/error"
        }
    }

    const createCode = () => {
        try{
            return true
        } catch(error) {
            return false
        }
    }

    const verifyCode = () => {
        try{
            setAuth(true)
        } catch(error) {
            setAuth(false)
        }
    }

    const createNewWallet = () => {
    }

    useEffect(() => {
        setTimeout(()=> {
            setIsReady(true)
        }, 2000)
    },[])

    return (
      <WalletContext.Provider
        {...props}
        value={{
            isReady,
            auth,
            walletExist,
            createNewWallet,
            createCode,
            verifyCode,
        }}
      />
    );
  };
  
  export {WalletProvider, useWallet};