import { Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';

export default function SignUpPage() {
  const { register, handleSubmit } = useForm();
  const navigate = useNavigate();

  const onSubmit = (data: any) => {
    console.log('Sign up data:', data);
    // TODO: Connect to backend IAM service
    navigate('/dashboard');
  };

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <h1 className="text-3xl font-display mb-2">Create Workspace.</h1>
      <p className="text-white/60 mb-8">Register your organization to EPMP.</p>

      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
        <div className="flex flex-col gap-2">
          <label className="text-sm text-white/80">Organization Name</label>
          <input 
            {...register('orgName')}
            type="text" 
            placeholder="Acme Properties Ltd."
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>

        <div className="flex flex-col gap-2">
          <label className="text-sm text-white/80">Admin Email</label>
          <input 
            {...register('email')}
            type="email" 
            placeholder="admin@acme.com"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>
        
        <div className="flex flex-col gap-2">
          <label className="text-sm text-white/80">Password</label>
          <input 
            {...register('password')}
            type="password" 
            placeholder="••••••••"
            className="bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white placeholder:text-white/30 focus:outline-none focus:border-orange focus:ring-1 focus:ring-orange transition-all"
            required
          />
        </div>

        <button type="submit" className="mt-4 bg-white text-black font-bold py-3 rounded-xl hover:bg-orange transition-colors">
          Register Organization
        </button>
      </form>

      <p className="text-center text-sm text-white/60 mt-8">
        Already registered?{' '}
        <Link to="/auth/signin" className="text-white hover:text-orange font-semibold transition-colors">Sign In</Link>
      </p>
    </div>
  );
}
