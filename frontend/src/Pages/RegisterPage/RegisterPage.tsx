import React, { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import checker from './registerController'
import './registerPage.css'
import isAuth from '../../utils/authCheck'

interface RegistrationCredentials {
    firstname: string;
    lastname: string;
    email: string;
    password: string;
    date: string;
    nickName?: string;
    avatar?: File;
    about?: string;
}

export interface RegisterPageProps {
    setLogin: React.Dispatch<React.SetStateAction<boolean>>;
}

const RegisterPage: React.FC<RegisterPageProps> = ({ setLogin }) => {
    const navigate = useNavigate()
    const [step, setStep] = useState(1)
    const [error, setError] = useState('')
    const [formData, setFormData] = useState<RegistrationCredentials>({
        firstname: '',
        lastname: '',
        email: '',
        password: '',
        date: '',
    });

    useEffect(() => {
        isAuth().then((isLogged) => {
            if (isLogged) {
                navigate('/')
            }
        })
    }, [])

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value, type } = e.target
        setFormData((prevData) => ({
            ...prevData,
            [name]: type === 'file' ? e.target.files?.[0] : value
        }))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (step === 1) {
            checker(formData) ? setError(checker(formData)) : setStep(2)
        } else {
            //Deconstruct the obj, to get the correct format of JSON
            const { avatar, ...formDataWithoutAvatar } = formData

            await fetch('http://localhost:8080/api/register', {
                method: "POST",
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(formDataWithoutAvatar),
                credentials: "include"
            })
                .then(async (response) => { //if we recieve error from backend
                    if (!response.ok) {
                        const message = await response.text()
                        setError(message) //show the message of error
                        return
                    } else if (response.ok) {
                        if (!avatar) {
                            setLogin(true)
                            navigate('/')
                        } else {
                            const data = new FormData() //FormData is for Image to send this as File 
                            data.append('image', avatar)
                            await fetch('http://localhost:8080/api/image/user', {
                                method: "POST",
                                body: data,
                                credentials: "include"
                            }).then(() => {
                                setLogin(true)
                                navigate('/')
                            })
                        }
                    }
                }).catch((error) => { //catch error if there is problem with fetch method
                    console.error(error)
                })
        }
    }

    return (
        <div className="loginpage">
            <div className='login-page__wrapper'>
                <div className="welcome-back-container">
                    <span className="welcome-back-glow">
                        <p className="welcome">WELCOME</p>
                        <p className="welcome">TO</p>
                        <p className="welcome">S0c1al-</p>
                        <p className="welcome">N3TW0rK</p>
                    </span>
                </div>
                <div className="register-page__form--wrapper">
                    <div className="loginTitle register__title">REGISTER</div>

                    <form className="register-page__form" onSubmit={handleSubmit}>
                        {step === 1 && (
                            <>
                                <input
                                    className="loginInput register__input"
                                    placeholder="firstname"
                                    name="firstname"
                                    value={formData.firstname}
                                    onChange={handleChange}
                                />
                                <input
                                    className="loginInput register__input"
                                    placeholder="lastname"
                                    name="lastname"
                                    value={formData.lastname}
                                    onChange={handleChange}
                                />
                                <input
                                    className="loginInput register__input"
                                    placeholder="email"
                                    type='email'
                                    name="email"
                                    value={formData.email}
                                    onChange={handleChange}
                                />
                                <input
                                    className="loginInput register__input"
                                    placeholder="password"
                                    type="password"
                                    name="password"
                                    value={formData.password}
                                    onChange={handleChange}
                                />
                                <input
                                    className="loginInput register__input"
                                    type="date"
                                    name="date"
                                    value={formData.date}
                                    onChange={handleChange}
                                    placeholder="date of birth"
                                />
                                <button className="button register__button" type='submit'>NEXT</button>
                            </>
                        )}

                        {step === 2 && (
                            <>
                                <input
                                    className="loginInput register__input"
                                    type="file"
                                    name="avatar"
                                    onChange={handleChange}
                                />
                                <input
                                    className="loginInput register__input"
                                    placeholder="About me"
                                    name="about"
                                    value={formData.about}
                                    onChange={handleChange}
                                />
                                <input
                                    className="loginInput register__input"
                                    placeholder="Nickname"
                                    name="nickName"
                                    value={formData.nickName}
                                    onChange={handleChange}
                                />
                                <div className="buttons-div">
                                    <button className="register__button button"
                                        onClick={() => { setStep(1) }}>BACK</button>
                                    <button className="button register__button">SUBMIT</button>
                                </div>
                            </>
                        )}
                    </form>
                    {error !== '' &&
                        <p className='auth-error-text'>{error}</p>
                    }
                    <Link to="/login" className="link-to-login link-to-register">
                        &lt;Already have an account?&gt;
                    </Link>
                </div>
            </div>
        </div>
    )
}

export default RegisterPage
