function Main() {
    return (
        <div id="main">
            <div id="main-navigation">
                <button>
                    <h1>Send Payment</h1>
                </button>
            </div>

            <div className='row'>
                <div id="balance-card">
                    <div>
                        <p>Total Balance</p>
                    </div>
                    <div>
                        <h1>$ 0</h1>
                    </div>
                </div>

                <div id="details-card">
                    <div>
                        <h1>Account Details</h1>
                        <p>Address: </p>
                        <p>Network: </p>
                    </div>
                </div>
            </div>
            <div className='row'>
                <div id="transaction-card">
                    <div>
                        <h1>Transactions</h1>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default Main
