import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';

export default function SignInPage() {
  const { register, handleSubmit } = useForm();
  const navigate = useNavigate();

  const onSubmit = (data: any) => {
    console.log('Sign in data:', data);
    // TODO: Connect to backend IAM service
    // Untuk sementara, langsung redirect ke dashboard agar mempermudah testing
    navigate('/dashboard');
  };

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <h1 className="text-3xl font-display mb-2">Welcome back.</h1>
      <p className="text-white/60 mb-8">Sign in to your EPMP enterprise account.</p>

      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
        <div className="flex flex-col gap-2">
          <label className="text-sm text-white/80">Email Address</label>
          <input 
            {...register('email')}
            type="email" 
            placeholder="admin@enterprise.com"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>
        
        <div className="flex flex-col gap-2">
          <div className="flex justify-between items-center">
            <label className="text-sm text-white/80">Password</label>
            <Link to="#" className="text-xs text-orange hover:underline">Forgot password?</Link>
          </div>
          <input 
            {...register('password')}
            type="password" 
            placeholder="••••••••"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>

        <button type="submit" className="mt-4 bg-orange text-black font-bold py-3 rounded-xl hover:bg-white transition-colors">
          Sign In
        </button>
      </form>

      <p className="text-center text-sm text-white/60 mt-8">
        Don't have an organization account?{' '}
        <Link to="/auth/signup" className="text-white hover:text-orange font-semibold transition-colors">Request Access</Link>
      </p>
    </div>
  );
}
