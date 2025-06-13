import React, { Dispatch, RefObject, SetStateAction, useState } from 'react';
import './emojiPicker.css'

const emojis = [
    {
        unicode: '\u{1F603}',
        name: 'SMILE'
    },
    {
        unicode: '\u{1F604}',
        name: 'Ha-Ha'
    },
    {
        unicode: '\u{1F600}',
        name: 'He-he'
    },
    {
        unicode: '\u{1F601}',
        name: 'XaXaxa'
    },
    {
        unicode: '\u{1F606}',
        name: 'XD'
    },
    {
        unicode: '\u{1F605}',
        name: 'hahaaaa'
    },
    {
        unicode: '\u{1F923}',
        name: 'LMFAO'
    },
    {
        unicode: '\u{1F602}',
        name: 'I`m crying'
    },
    {
        unicode: '\u{1F642}',
        name: ')'
    },
    {
        unicode: '\u{1F643}',
        name: 'What??))'
    },
    {
        unicode: '\u{1F609}',
        name: ';)'
    },
    {
        unicode: '\u{1F60A}',
        name: 'Nyam'
    },
    {
        unicode: '\u{1F607}',
        name: 'Saint'
    },
    {
        unicode: '\u{1F61B}',
        name: 'MLEM'
    },
    {
        unicode: '\u{1F61C}',
        name: 'Blink Mem'
    },
    {
        unicode: '\u{1F60B}',
        name: 'Yammy'
    },
    {
        unicode: '\u{1F92A}',
        name: 'Crazy'
    },
    {
        unicode: '\u{1F61D}',
        name: 'Crazy Laugh'
    },
    {
        unicode: '\u{1F911}',
        name: 'MONEYYYY'
    },
    {
        unicode: '\u{1F636}',
        name: 'NO WORDS'
    },
    {
        unicode: '\u{1F62A}',
        name: 'Zzzzz...'
    },
    {
        unicode: '\u{1F922}',
        name: 'DiSGASTIENG!'
    },
    {
        unicode: '\u{1F60E}',
        name: 'Cool asf'
    },
    {
        unicode: '\u{1F61F}',
        name: 'What? :((('
    },
    {
        unicode: '\u{1F641}',
        name: ':('
    },
    {
        unicode: '\u{1F620}',
        name: '><'
    },
    {
        unicode: '\u{1F92C}',
        name: '#$!&*!()!@)_'
    },
    {
        unicode: '\u{1F480}',
        name: 'I`m dead'
    },
    {
        unicode: '\u{1F4A9}',
        name: 'SHIT! POOP!'
    },
    {
        unicode: '\u{1F921}',
        name: 'I`m HTML developer'
    },
    {
        unicode: '\u{1F47D}',
        name: 'WHO IS IN PARIS?!?!'
    },
    {
        unicode: '\u{1F49A}',
        name: 'LoVe'
    },
    {
        unicode: '\u{1F44C}',
        name: 'OK'
    },
    {
        unicode: '\u{1F918}',
        name: 'ROCK'
    },
    {
        unicode: '\u{1F44F}',
        name: '*Clap*'
    },
    {
        unicode: '\u{1F90F}',
        name: 'Lil` bit'
    },
    {
        unicode: '\u{1F595}',
        name: 'F U'
    },
    {
        unicode: '\u{1F44D}',
        name: 'Bread'
    }
];

interface EmojiPickerProps {
    input: RefObject<HTMLTextAreaElement> | RefObject<HTMLInputElement>
    closeBtn: (event: React.MouseEvent<HTMLButtonElement>) => void
    className: string
    setData?: Dispatch<SetStateAction<{ chatId: number; content: string; type: string; }>>
}

const EmojiPicker: React.FC<EmojiPickerProps> = ({ input, closeBtn, className, setData }) => {

    const addEmoji: React.MouseEventHandler<HTMLSpanElement> = (event) => {
        const emoji = event.target as HTMLElement
        const textarea = input.current
        if (textarea) {
            const cursorPos = textarea.selectionStart || 0
            const textBeforeCursor = textarea.value.substring(0, cursorPos)
            const textAfterCursor = textarea.value.substring(cursorPos)
            const updatedText = textBeforeCursor + emoji.innerText + textAfterCursor

            textarea.value = updatedText
            textarea.focus()

            setData && setData((prev) => {
                return {
                    ...prev,
                    content: updatedText
                }
            })

            const inputEvent = new Event('input', { bubbles: true })
            textarea.dispatchEvent(inputEvent)
        }
    }

    return (
        <div className={`emoji-list--container ${className}`} >
            <div className='emoji-list'>
                {emojis.map((emoji, key) => (
                    <span className='emoji' data-name={emoji.name}
                        id={emoji.unicode} onClick={addEmoji} key={key}>{emoji.unicode}</span>
                ))}
            </div>
            <button className='emoji-list__close' onClick={closeBtn}>X</button>
        </div>
    )
}

export default EmojiPicker