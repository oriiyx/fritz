export function Footer() {
    const currentYear = new Date().getFullYear()

    return (
        <footer className="footer footer-center bg-base-100 p-4 text-base-content shadow-inner">
            <aside>
                <p>
                    Copyright © {currentYear} - F3D
                </p>
            </aside>
        </footer>
    )
}
