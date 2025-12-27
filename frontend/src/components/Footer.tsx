import { Link } from 'react-router-dom'
import './Footer.css'

const Footer = () => {
  const currentYear = new Date().getFullYear()

  return (
    <footer className="app-footer">
      <div className="footer-content">
        <div className="footer-links">
          <Link to="/privacy-policy" className="footer-link">Privacy Policy</Link>
          <span className="footer-separator">•</span>
          <Link to="/terms-of-service" className="footer-link">Terms of Service</Link>
        </div>
        <div className="footer-copyright">
          © {currentYear} Social Media Publisher. All rights reserved.
        </div>
      </div>
    </footer>
  )
}

export default Footer
