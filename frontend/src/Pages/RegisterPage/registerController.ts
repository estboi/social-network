export default function checker(data: {
    firstname: string
    lastname: string
    email: string
    password: string
    date?: string
    avatar?: File
    about?: string
}): string {
    switch (true) {
        case !data.firstname:
            return 'First name is required.'
        case !data.lastname:
            return 'Last name is required.'
        case !data.email:
            return 'Email is required.'
        case !isValidEmail(data.email):
            return 'Invalid email format.'
        case !data.password:
            return 'Password is required.'
        case data.password.length < 6:
            return 'Password must be at least 6 characters long.'
        default:
            return ''
    }
}

function isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
}
