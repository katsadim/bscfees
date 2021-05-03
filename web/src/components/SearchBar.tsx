import * as React from "react";
import {FunctionComponent, useState} from "react";
import {IconButton, InputAdornment, TextField} from "@material-ui/core";
import AccountBalanceWalletOutlinedIcon from '@material-ui/icons/AccountBalanceWalletOutlined';
import SearchIcon from '@material-ui/icons/Search';

type Props = {
    handleSearch: (searchTerm: string) => void
    maxLength: number
}


export const SearchBar: FunctionComponent<Props> = (props: Props) => {
    const [error, setError] = useState(false);
    const [disabledSearch, setDisabledSearch] = useState(true);
    const [prevValue, setPrevValue] = useState('');
    const [currentValue, setCurrentValue] = useState('');

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>): void => {
        setPrevValue(event.target.value);
    };

    const validateInput = (event: React.ChangeEvent<HTMLInputElement>): void => {
        if (event.target.value === "") {
            setDisabledSearch(true)
            return
        }
        const lastLetter = event.target.value[event.target.value.length - 1].toLowerCase()
        if (((lastLetter >= 'a' && lastLetter <= 'z')
            || ((lastLetter >= '0' && lastLetter <= '9')))
            && event.target.value.length <= 42) {
            setError(false)
            setDisabledSearch(false)
            setCurrentValue(event.target.value)
        } else {
            event.target.value = prevValue
            setError(true)
            setDisabledSearch(true)
        }
    }

    const performSearch = (): void => {
        props.handleSearch(currentValue)
    }

    return (
        <TextField
            id="bsc-address"
            name="bsc-address"
            aria-label="wallet address"
            label="Wallet Address"
            variant="outlined"
            defaultValue="0x"
            onChange={handleChange}
            onInput={validateInput}
            error={error}
            helperText={error ? `Should be a ${props.maxLength} long character hexadecimal` : ''}
            InputProps={{
                startAdornment: (
                    <InputAdornment position="start">
                        <IconButton
                            aria-label="wallet"
                            disabled
                        >
                            <AccountBalanceWalletOutlinedIcon/>
                        </IconButton>


                    </InputAdornment>
                ),
                endAdornment: (
                    <InputAdornment position="end">
                        <IconButton
                            aria-label="search"
                            disabled={disabledSearch}
                            type="submit"
                            onClick={performSearch}
                        >
                            <SearchIcon/>
                        </IconButton>
                    </InputAdornment>),
            }}
            fullWidth
        />
    );
};

export default SearchBar;
