import React, { useEffect, useState } from 'react'
import './login.css'
import { Link, useNavigate } from 'react-router-dom'
import isAuth from '../../utils/authCheck'
import { RegisterPageProps } from '../RegisterPage/RegisterPage'



const LoginPage: React.FC<RegisterPageProps> = ({ setLogin }) => {
    const navigate = useNavigate()
    const [formData, setFormData] = useState({ login: '', password: '' })
    const [error, setError] = useState('')

    /** Validate user input */
    const ValidData = (data: { login: string, password: string }) => { //Validate data on submit 
        if (data.login === '') return 'Email/UserName is required'
        //regex test if the string has @ -> it needs to be in correct format
        if (data.login.includes('@') && !/^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]+/.test(data.login)) {
            return 'Email should be in correct format'
        }
        if (data.password === '') { return 'Password is required' }
        if (data.password.length < 6) { return 'Password is minimum 6 characters' }
        return ''
        
    }

    useEffect(() => {
        isAuth().then((isLogged) => {
            if (isLogged === true) {
                navigate('/')
            }
        })
    }, [])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        let message: string
        if (message = ValidData(formData)) {
            setError(message)
            return
        }
        await fetch('http://localhost:8080/api/login', {
            method: "POST",
            headers: { "Content-Type": "application/json", },
            credentials: "include",
            body: JSON.stringify(formData)
        })
            .then(async (response) => { //if we recieve error from backend
                switch (response.status) {
                    case 200:
                        console.log('Logged succesfully')
                        setLogin(true)
                        navigate('/')
                        break
                    case 404:
                        console.log('NOT FOUND')
                        break
                    default:
                        const errMsg = await response.text()
                        setError(errMsg)
                }
            }).catch((error) => { //catch error if there is problem with fetch method
                console.error(error)
            })
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target
        setFormData((prevData) => ({
            ...prevData,
            [name]: value,
        }))
    }

    return (
        <div className="loginpage">
            <div className='login-page__wrapper'>
                <div className="welcome-back-container">
                    <span className="welcome-back-glow">
                        <p className="welcome">WELCOME</p>
                        <p className="welcome">BACK TO</p>
                        <p className="welcome">S0c1al-</p>
                        <p className="welcome">N3TW0rK</p>
                    </span>
                </div>
                <div className="form-parent">
                    <div className="loginTitle">LOGIN</div>
                    <form className="loginForm" onSubmit={handleSubmit}>
                        <input
                            className="loginInput"
                            placeholder="email || username"
                            name="login"
                            type='text'
                            maxLength={254}
                            value={formData.login}
                            onChange={handleChange}

                        />
                        <input
                            className="loginInput"
                            placeholder="password"
                            type="password"
                            name="password"
                            maxLength={254}
                            value={formData.password}
                            onChange={handleChange}
                        />
                        <button className="button">SUBMIT</button>
                    </form>
                    {error != '' && <p className='auth-error-text'>{error}</p>}
                    <Link to="/register" className="link-to-register">
                        &lt;Don't have account?&gt;
                    </Link>
                </div>
            </div>
        </div>
    )
}
export default LoginPage
