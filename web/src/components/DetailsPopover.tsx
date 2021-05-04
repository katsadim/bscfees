import * as React from "react";
import {FunctionComponent} from "react";
import {createStyles, IconButton, makeStyles, Popover, Theme, Typography} from "@material-ui/core";
import InfoOutlinedIcon from "@material-ui/icons/InfoOutlined";


type Props = {
    bnbusdPrice: number
    ethusdPrice: number
}

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        typography: {
            padding: theme.spacing(1),
        },
    }),
);

export const DetailsPopover: FunctionComponent<Props> = (props: Props) => {
    const classes = useStyles();
    const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement | null>(null);

    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const open = Boolean(anchorEl);
    const id = open ? 'simple-popover' : undefined;

    return (
        <div>
            <IconButton aria-label="share" aria-describedby={id} onClick={handleClick}>
                <InfoOutlinedIcon/>
            </IconButton>
            <Popover
                id={id}
                open={open}
                anchorEl={anchorEl}
                onClose={handleClose}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'center',
                }}
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'center',
                }}
            >
                <Typography className={classes.typography}>Current BNB/BUSD price: {props.bnbusdPrice}</Typography>
                <Typography className={classes.typography}>Current ETH/BUSD price: {props.ethusdPrice}</Typography>
            </Popover>
        </div>

    )
}
