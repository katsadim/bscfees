import * as React from "react"
import {FunctionComponent, useState} from "react"
import {
    createMuiTheme,
    CssBaseline,
    Grid,
    IconButton,
    LinearProgress,
    makeStyles,
    Snackbar,
    Theme,
    ThemeProvider,
    Typography,
    useMediaQuery
} from "@material-ui/core";
import CloseIcon from '@material-ui/icons/Close';
import SearchBar from "../components/SearchBar";
import {Fees} from "../http/Response";
import {api} from "../http/Api";
import {RequestErrorType} from "../model/Error";
import {ResultCard} from "../components/ResultCard";
import Layout from "../components/Layout";
import {SEO} from "../components/Seo";

const useStyles = makeStyles((theme) => ({
    linearProgress: {
        width: '100%',
        '& > * + *': {
            marginTop: theme.spacing(2),
        },

    },
    container: {
        maxWidth: "100%"
    },
}));

const Index: FunctionComponent = () => {
    const classes = useStyles();
    const [searchInitiated, setSearchInitiated] = useState(false);
    const [searchSuccess, setSearchSuccess] = useState(false);
    const [openSnackbar, setOpenSnackbar] = React.useState(false);
    const [errorMessage, setErrorMessage] = React.useState("");
    const [fees, setFees] = React.useState(0.0);
    const [bnbusdPrice, setBnbUsdPrice] = React.useState(0.0);

    const performSearch = (searchTerm: string): void => {
        const error = validateSearchTerm(searchTerm);
        if (error != RequestErrorType.OK) {
            setErrorMessage(error)
            setOpenSnackbar(true)
            return
        }

        setSearchInitiated(true)
        setSearchSuccess(false)
        api<Fees>(`fees?address=${searchTerm}`)
            .then((f) => {
                setFees(f.fees)
                setBnbUsdPrice(f.bnbusdPrice)
                setSearchSuccess(true)
            })
            .catch(error => {
                if (error == "TypeError: Failed to fetch") {
                    throw Error('Network Connection Error')
                }
                setErrorMessage(RequestErrorType.WALLET_NOT_FOUND)
                setOpenSnackbar(true)
            })
            .catch(() => {
                setErrorMessage(RequestErrorType.NETWORK_ERROR)
                setOpenSnackbar(true)
            })
            .finally(() => setSearchInitiated(false))
    }

    const handleClose = (event: React.SyntheticEvent | React.MouseEvent, reason?: string) => {
        if (reason === 'clickaway') {
            return;
        }

        setOpenSnackbar(false);
    };

    const theme = determineTheme()

    return (
        <React.StrictMode>
            <SEO/>
            <ThemeProvider theme={theme}>
                <CssBaseline/>
                <Layout>
                    <Grid container spacing={4} justifyContent="center" className={classes.linearProgress}>
                        <Grid container item xs={12} justifyContent="center">
                            <Typography variant="h4" component="h1" gutterBottom>
                                BscFees
                            </Typography>
                        </Grid>

                        <Grid container item spacing={0} justifyContent="center">
                            <Grid container item xs={12} justifyContent="center">
                                <SearchBar
                                    handleSearch={searchTerm => performSearch(searchTerm)}
                                    maxLength={42}
                                />
                            </Grid>
                            <Grid container item xs={12} justifyContent="center">
                                <div className={classes.linearProgress}>
                                    {searchInitiated ? <LinearProgress/> : <></>}
                                </div>
                            </Grid>
                        </Grid>

                        <Grid container item xs={12} justifyContent="center" zeroMinWidth>
                            {searchSuccess ? <ResultCard bnbusdPrice={bnbusdPrice} fees={fees} numOfDecimals={4}/> : <></>}
                        </Grid>
                    </Grid>
                </Layout>
                <Snackbar
                    anchorOrigin={{
                        vertical: 'bottom',
                        horizontal: 'left',
                    }}
                    open={openSnackbar}
                    autoHideDuration={6000}
                    onClose={handleClose}
                    message={errorMessage}
                    action={
                        <React.Fragment>
                            <IconButton size="small" aria-label="close" color="inherit" onClick={handleClose}>
                                <CloseIcon fontSize="small"/>
                            </IconButton>
                        </React.Fragment>
                    }
                />
            </ThemeProvider>
        </React.StrictMode>
    );
};

const validateSearchTerm = (searchTerm: string): RequestErrorType => {
    if (searchTerm[0] != '0' && searchTerm[1].toLowerCase() != 'x') {
        return RequestErrorType.NOT_STARTING_WITH_0X
    } else if (searchTerm.length != 42) {
        return RequestErrorType.NO_APPROPRIATE_LENGTH
    }
    return RequestErrorType.OK
}

const determineTheme = (): Theme => {
    const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');

    return React.useMemo(
        () =>
            createMuiTheme({
                palette: {
                    mode: prefersDarkMode ? 'dark' : 'light',
                },
            }),
        [prefersDarkMode],
    );
}

export default Index;