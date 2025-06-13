import React, { useState } from "react";
import "./searchBar.css";

interface Props {
    placeHolder: string;
    type: string;
}

const SearchBar = ({ placeHolder, type }: Props) => {
    const [inputValue, changeValue] = useState('')
    const [previousChoose, setPrevious] = useState<HTMLElement>()

    const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        const elements = document.querySelectorAll(`[data-type=${type}]`);
        elements.forEach(element => {
            if (element.textContent && element.textContent === inputValue) {
                const parent = element.parentElement
                if (parent) {
                    if (previousChoose) previousChoose.dataset.choosed = 'false'
                    setPrevious(parent)
                    parent.dataset.choosed = 'true'
                    parent.scrollIntoView({ behavior: 'smooth', block: 'start', inline: 'nearest' });
                }
            }
        });
    }
    const handleChange = (event: React.FormEvent<HTMLInputElement>) => {
        event.preventDefault()
        changeValue(event.currentTarget?.value)
    }
    return (
        <form className="search-bar--wrapper" onSubmit={handleSubmit}>
            <input list='search__datalist' className="search-bar" placeholder={placeHolder} onChange={handleChange} />
            <button type="submit" className="search-bar__submit" >
                <img alt='submit' src="/assets/Submit.svg"/>
            </button>
            <datalist id='search__datalist' />
        </form>
    );
};
export default SearchBar