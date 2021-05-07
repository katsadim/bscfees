import * as React from "react";
import {FunctionComponent} from "react";
import {Box, Container, Grid, IconButton, Tooltip, Typography} from "@material-ui/core";
import Link from '@material-ui/core/Link';
import {BscScanIcon} from "./icons/BscScanIcon";
import GitHubIcon from '@material-ui/icons/GitHub';
import "../index.css";
import {OutboundLink} from "gatsby-plugin-google-gtag";
import {EthScanIcon} from "./icons/EthScanIcon";

export interface LayoutProps {
    children: React.ReactNode
}


const Copyright: FunctionComponent = () => {

    return (
        <Grid container justifyContent="center">
            <Grid item xs={12} justifyContent="center">
                <Typography variant="body2" color="text.secondary" align="center">
                    {'Copyright © '}
                    <Link color="inherit" href="https://www.bscfees.com/">
                        BscFees
                    </Link>{' '}
                    {new Date().getFullYear()}
                    {'.'}
                </Typography>
            </Grid>
            <Grid item xs={12} justifyContent="center">
                <Typography variant="body2" color="text.secondary" align="center">
                    {'Powered by:'}
                </Typography>
            </Grid>
            <Grid container item justifyContent="center" alignItems={"flex-start"} alignContent={"center"}>
                <Grid item/>
                <Grid item>
                    <Tooltip title="EtherScan">
                        <OutboundLink href="https://www.etherscan.io">
                            <IconButton aria-label="etherscan" onClick={() => window.open('https://www.etherscan.io')}>
                                <EthScanIcon width={20} height={20} style={{fontSize: 10}}/>
                            </IconButton>
                        </OutboundLink>
                    </Tooltip>
                </Grid>
                <Grid item>
                    <Tooltip title="BscScan">
                        <OutboundLink href="https://www.bscscan.com">
                            <IconButton aria-label="bscscan" onClick={() => window.open('https://www.bscscan.com')}>
                                <BscScanIcon width={20} height={20} style={{fontSize: 10}}/>
                            </IconButton>
                        </OutboundLink>
                    </Tooltip>
                </Grid>
                <Grid item>
                    <Tooltip title="Github">
                        <OutboundLink href="https://www.github.com/katsadim/bscfees">
                            <IconButton aria-label="github"
                                        onClick={() => window.open('https://www.github.com/katsadim/bscfees')}>
                                <GitHubIcon style={{fontSize: 20}}/>
                            </IconButton>
                        </OutboundLink>
                    </Tooltip>
                </Grid>
                <Grid item/>
            </Grid>
        </Grid>
    );
}

const Layout: FunctionComponent<LayoutProps> = (props: LayoutProps) => {
    return (
        <React.StrictMode>
            <div style={{
                height: `100%`,
                display: `flex`,
                flexDirection: `column`,
            }}>
                <Container maxWidth="sm" style={{flex: 1}}>
                    <Box my={4}>
                        {props.children}
                    </Box>
                </Container>
                <Copyright/>
            </div>
        </React.StrictMode>
    )
}

export default Layout;