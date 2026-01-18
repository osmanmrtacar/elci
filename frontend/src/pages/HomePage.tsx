import {
  Header,
  HeroSection,
  HowItWorks,
  Benefits,
  WhoItsFor,
  TrustSection,
  FinalCTA,
  Footer
} from '../components/landing'

const HomePage = () => {
  return (
    <div className="min-h-screen bg-white">
      <Header />
      <main>
        <HeroSection />
        <HowItWorks />
        <Benefits />
        <WhoItsFor />
        <TrustSection />
        <FinalCTA />
      </main>
      <Footer />
    </div>
  )
}

export default HomePage
