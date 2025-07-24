import { Link } from "react-router-dom";
import Header from "../components/Header";

export default function PrivacyPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-sky-100 to-teal-50">
      <Header />
      
      <main className="max-w-4xl mx-auto px-4 py-12">
        <div className="bg-white/80 rounded-2xl shadow-xl backdrop-blur p-8 md:p-12">
          {/* Header Section */}
          <div className="text-center mb-12">
            <h1 className="text-4xl font-bold text-gray-900 mb-4">
              Privacy Policy
            </h1>
            <p className="text-lg text-gray-600">
              Last updated: {new Date().toLocaleDateString()}
            </p>
          </div>

          {/* Content */}
          <div className="space-y-8 text-gray-700 leading-relaxed">
            
            <div className="bg-red-50 border-l-4 border-red-400 p-6 rounded-r-lg">
              <div className="flex">
                <div className="flex-shrink-0">
                  <svg className="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                  </svg>
                </div>
                <div className="ml-3">
                  <h3 className="text-sm font-medium text-red-800">
                    IMPORTANT DISCLAIMER
                  </h3>
                  <div className="mt-2 text-sm text-red-700">
                    <p>
                      BY USING THIS SERVICE, YOU ACKNOWLEDGE AND AGREE THAT WE HAVE ZERO LIABILITY 
                      FOR ANY CONTENT, DATA LOSS, PRIVACY BREACHES, OR CONSEQUENCES ARISING FROM YOUR USE OF THIS PLATFORM.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Data Collection</h2>
              <p className="mb-4">
                We collect minimal information to provide our service:
              </p>
              <ul className="list-disc ml-6 space-y-2">
                <li>Email addresses and usernames for registered accounts</li>
                <li>Messages sent in chat rooms (temporarily stored)</li>
                <li>Basic usage analytics and system logs</li>
                <li>Cookies for authentication and functionality</li>
              </ul>
              <p className="mt-4 font-semibold text-red-600">
                WARNING: We cannot guarantee the security of any data you provide. Use this service at your own risk.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Data Retention & Deletion</h2>
              <p className="mb-4">
                Our platform operates on a 24-hour cycle:
              </p>
              <ul className="list-disc ml-6 space-y-2">
                <li>Chat messages are deleted when rooms expire (24 hours)</li>
                <li>Account information is retained until account deletion</li>
                <li>System logs may be retained for operational purposes</li>
                <li>We reserve the right to retain or delete any data without notice</li>
              </ul>
              <p className="mt-4 font-semibold text-red-600">
                IMPORTANT: Data deletion is not guaranteed. Technical failures may result in data persistence beyond intended timeframes.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Third-Party Services</h2>
              <p className="mb-4">
                Our service may integrate with third-party APIs and services for content aggregation.
                We have no control over their privacy practices.
              </p>
              <p className="font-semibold text-red-600">
                We are not responsible for any third-party data handling, breaches, or privacy violations.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">User Responsibility</h2>
              <p className="mb-4">
                You are solely responsible for:
              </p>
              <ul className="list-disc ml-6 space-y-2">
                <li>All content you post or share</li>
                <li>Protecting your account credentials</li>
                <li>Any consequences of your participation</li>
                <li>Verifying the privacy and security of your communications</li>
                <li>Understanding that anonymous does not mean secure or private</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Security Limitations</h2>
              <p className="mb-4 font-semibold text-red-600">
                WE PROVIDE NO GUARANTEES REGARDING:
              </p>
              <ul className="list-disc ml-6 space-y-2 text-red-600">
                <li>Data encryption or security</li>
                <li>Protection from data breaches</li>
                <li>Anonymity or privacy preservation</li>
                <li>Uptime or availability</li>
                <li>Data backup or recovery</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">No Liability</h2>
              <div className="bg-gray-100 p-6 rounded-lg">
                <p className="font-bold text-gray-900 mb-4">
                  COMPLETE DISCLAIMER OF LIABILITY:
                </p>
                <p className="text-sm">
                  WE DISCLAIM ALL LIABILITY FOR ANY DAMAGES, LOSSES, OR CONSEQUENCES 
                  ARISING FROM YOUR USE OF THIS SERVICE, INCLUDING BUT NOT LIMITED TO:
                  DATA LOSS, PRIVACY BREACHES, IDENTITY EXPOSURE, HARASSMENT, TRAUMA, 
                  FINANCIAL LOSS, EMOTIONAL DISTRESS, OR ANY OTHER HARM. YOUR USE OF 
                  THIS SERVICE IS ENTIRELY AT YOUR OWN RISK.
                </p>
              </div>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Changes to Privacy Policy</h2>
              <p>
                We may update this privacy policy at any time without notice. 
                Continued use of the service constitutes acceptance of any changes.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Contact</h2>
              <p>
                For privacy-related questions, you may contact us through the platform.
                However, we make no guarantees regarding response times or resolution of concerns.
              </p>
            </section>
          </div>

          {/* Back Link */}
          <div className="mt-12 text-center">
            <Link
              to="/about"
              className="text-indigo-600 hover:text-indigo-500 transition-colors"
            >
              ← Back to About
            </Link>
          </div>
        </div>
      </main>
    </div>
  );
}