import React, {
    createContext,
    ReactNode,
    useContext,
    useEffect,
  } from 'react';

  
  interface Wallet {
    isReady: boolean;
    auth: boolean;
    walletExist: () => boolean,
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
    const [isReady, setIsReady] = React.useState(false)
    const [auth, setAuth] = React.useState(false)

    const walletExist = () => {
        try{
            return false
        } catch(error) {
            return false
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
        })
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