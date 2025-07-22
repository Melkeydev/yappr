import { Link } from "react-router-dom";
import Header from "../components/Header";

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-50 via-sky-100 to-teal-50">
      <Header />
      
      <main className="max-w-4xl mx-auto px-4 py-12">
        <div className="bg-white/80 rounded-2xl shadow-xl backdrop-blur p-8 md:p-12">
          {/* Header Section */}
          <div className="text-center mb-12">
            <h1 className="text-4xl font-bold text-gray-900 mb-4">
              Terms of Service
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
                    CRITICAL WARNING
                  </h3>
                  <div className="mt-2 text-sm text-red-700">
                    <p>
                      BY USING THIS SERVICE, YOU ACCEPT FULL RESPONSIBILITY FOR ALL RISKS AND 
                      ACKNOWLEDGE OUR COMPLETE DISCLAIMER OF LIABILITY FOR ANY AND ALL CONSEQUENCES.
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Acceptance of Terms</h2>
              <p>
                By accessing or using Anonymous Chat, you agree to be bound by these Terms of Service.
                If you disagree with any part of these terms, you must not use our service.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Service Description</h2>
              <p className="mb-4">
                Anonymous Chat is a platform for temporary, anonymous conversations. We provide:
              </p>
              <ul className="list-disc ml-6 space-y-2">
                <li>24-hour chat rooms with daily topic rotations</li>
                <li>Anonymous messaging capabilities</li>
                <li>Basic user account management</li>
              </ul>
              <p className="mt-4 font-semibold text-red-600">
                WE PROVIDE NO GUARANTEES REGARDING SERVICE AVAILABILITY, FUNCTIONALITY, OR SECURITY.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">User Responsibilities and Conduct</h2>
              <p className="mb-4">You are solely responsible for:</p>
              <ul className="list-disc ml-6 space-y-2">
                <li>All content you post, share, or communicate</li>
                <li>Your interactions with other users</li>
                <li>Any consequences arising from your use of the service</li>
                <li>Maintaining the security of your account</li>
                <li>Complying with applicable laws and regulations</li>
              </ul>
              <div className="bg-yellow-50 border border-yellow-200 p-4 rounded-lg mt-4">
                <p className="text-yellow-800 font-semibold">
                  WARNING: We do not monitor, moderate, or control user content or interactions.
                </p>
              </div>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Content and Communication</h2>
              <p className="mb-4 font-semibold text-red-600">
                WE ACCEPT NO RESPONSIBILITY FOR:
              </p>
              <ul className="list-disc ml-6 space-y-2 text-red-600">
                <li>Any content posted by users</li>
                <li>Harmful, offensive, illegal, or inappropriate content</li>
                <li>Harassment, bullying, or abusive behavior</li>
                <li>Misinformation, false claims, or misleading content</li>
                <li>Privacy violations or doxxing</li>
                <li>Discussion of sensitive, traumatic, or triggering topics</li>
                <li>Mental health impacts from user interactions</li>
              </ul>
              <p className="mt-4 text-red-600 font-semibold">
                Users engage with all content and other users entirely at their own risk.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">No Moderation Policy</h2>
              <div className="bg-red-100 p-6 rounded-lg">
                <h3 className="font-bold text-red-900 mb-3">EXPLICIT DISCLAIMER:</h3>
                <p className="text-red-800">
                  We provide NO content moderation, monitoring, or filtering. Users may encounter:
                </p>
                <ul className="list-disc ml-6 mt-2 space-y-1 text-red-800">
                  <li>Offensive, disturbing, or harmful content</li>
                  <li>Harassment or abusive behavior</li>
                  <li>Illegal or inappropriate material</li>
                  <li>Triggering or traumatic discussions</li>
                  <li>Misinformation or dangerous advice</li>
                </ul>
                <p className="mt-4 font-bold text-red-900">
                  WE ARE NOT LIABLE FOR ANY PSYCHOLOGICAL, EMOTIONAL, OR OTHER HARM RESULTING FROM EXPOSURE TO USER CONTENT.
                </p>
              </div>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Data and Privacy Risks</h2>
              <p className="mb-4 font-semibold text-red-600">
                WE MAKE NO GUARANTEES REGARDING:
              </p>
              <ul className="list-disc ml-6 space-y-2">
                <li>Data security or protection from breaches</li>
                <li>True anonymity or privacy preservation</li>
                <li>Data retention or deletion timelines</li>
                <li>Protection from identity exposure</li>
                <li>Backup or recovery of user data</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Service Availability</h2>
              <p>
                We provide this service "AS IS" without warranties of any kind. We reserve the right to:
              </p>
              <ul className="list-disc ml-6 space-y-2 mt-2">
                <li>Discontinue the service at any time without notice</li>
                <li>Modify features or functionality without warning</li>
                <li>Experience downtime, data loss, or technical failures</li>
                <li>Suspend or terminate accounts without cause</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Limitation of Liability</h2>
              <div className="bg-gray-900 text-white p-6 rounded-lg">
                <h3 className="font-bold mb-4 text-yellow-400">COMPLETE DISCLAIMER OF LIABILITY:</h3>
                <p className="text-sm leading-relaxed">
                  TO THE MAXIMUM EXTENT PERMITTED BY LAW, WE DISCLAIM ALL LIABILITY FOR ANY DAMAGES, 
                  LOSSES, HARM, OR CONSEQUENCES OF ANY KIND ARISING FROM YOUR USE OF THIS SERVICE, 
                  INCLUDING BUT NOT LIMITED TO:
                </p>
                <ul className="list-disc ml-6 mt-3 space-y-1 text-sm">
                  <li>DIRECT, INDIRECT, INCIDENTAL, SPECIAL, OR CONSEQUENTIAL DAMAGES</li>
                  <li>DATA LOSS, CORRUPTION, OR UNAUTHORIZED ACCESS</li>
                  <li>PRIVACY BREACHES OR IDENTITY EXPOSURE</li>
                  <li>EMOTIONAL DISTRESS, TRAUMA, OR MENTAL HEALTH IMPACTS</li>
                  <li>HARASSMENT, ABUSE, OR HARMFUL USER INTERACTIONS</li>
                  <li>FINANCIAL LOSSES OR BUSINESS INTERRUPTION</li>
                  <li>PHYSICAL HARM OR SAFETY RISKS</li>
                  <li>LEGAL CONSEQUENCES FROM USER ACTIVITIES</li>
                </ul>
                <p className="mt-4 font-bold text-yellow-400">
                  YOUR USE OF THIS SERVICE IS ENTIRELY AT YOUR OWN RISK.
                </p>
              </div>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Indemnification</h2>
              <p>
                You agree to indemnify and hold harmless the service operators from any claims, 
                damages, losses, or expenses arising from your use of the service or violation of these terms.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Governing Law</h2>
              <p>
                These terms are governed by applicable local laws. Any disputes will be resolved 
                in appropriate courts, though we disclaim responsibility for legal costs or outcomes.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Changes to Terms</h2>
              <p>
                We reserve the right to modify these terms at any time without notice. 
                Continued use of the service constitutes acceptance of any changes.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">Contact Information</h2>
              <p>
                For questions about these terms, you may contact us through the platform.
                We make no guarantees regarding response times or issue resolution.
              </p>
            </section>

            <div className="bg-orange-50 border border-orange-200 p-6 rounded-lg">
              <p className="text-orange-800 font-bold text-center">
                BY USING THIS SERVICE, YOU ACKNOWLEDGE THAT YOU HAVE READ, UNDERSTOOD, 
                AND AGREE TO THESE TERMS AND OUR COMPLETE DISCLAIMER OF LIABILITY.
              </p>
            </div>
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