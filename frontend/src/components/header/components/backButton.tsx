import { useNavigate } from "react-router-dom"

const HeaderBackButton = () => {
    const navigate = useNavigate()

    return <button className="header__back" onClick={() => { navigate(-1) }}>
        <img className="header__back-img" src="/assets/BackButton.svg" />BACK
    </button>
}

export default HeaderBackButton