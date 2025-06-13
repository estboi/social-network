import { useNavigate } from "react-router-dom";
import HomeHeader from "./components/homePageComp";
import HeaderBackButton from "./components/backButton";



import "./header.css"
import './headerAdap.css'
import GroupsHeader from "./components/groupPageComp";
import GroupProfileHeader from "./components/groupProfileComp";

interface Props {
    className: string,
}
type listComponentsType = {
    [key: string]: JSX.Element
}
const Childs: listComponentsType = {
    '/': <HomeHeader />,
    '/post/[0-9]+': <HeaderBackButton />,
    '/createPost(?:/(\\d+))?$': <HeaderBackButton />,
    '/groups/created': <GroupsHeader />,
    '/groups/[0-9]+': <GroupProfileHeader />,
    '/groups/createGroup': <HeaderBackButton />,
    '/groups/[0-9]+/members': <HeaderBackButton />,
    '/groups/[0-9]+/createEvent': <HeaderBackButton />,
    '/events/[0-9]+': <HeaderBackButton />
}

const Header = ({ className }: Props) => {
    const navigate = useNavigate()

    const path = window.location.pathname
    const matchedKey = Object.keys(Childs).find(key => {
        const regex = new RegExp(`^${key}$`);
        return regex.test(path);
    });

    // Render the matched component
    const matchedComponent = matchedKey ? Childs[matchedKey] : null;

    const isPageActive = (currentPage: string | undefined) => { return path === currentPage }

    const NavigateToPage = (event: React.MouseEvent<HTMLElement>) => {
        const element = event.target as HTMLElement
        const currentPage = element.dataset.setRoute
        if (isPageActive(currentPage)) return
        if (currentPage !== undefined) navigate(currentPage)
    }

    return (
        <div className={`header ${className}`}>
            <div className="header__title" data-set-route='/' onClick={NavigateToPage}>
                <p data-set-route='/' onClick={NavigateToPage}>S0c1al-N3TW0rK</p>
            </div>
            {matchedComponent}
        </div>
    )
}

export default Header