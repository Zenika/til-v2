<script lang="ts">
    let { children } = $props();
    import Header from "../components/header.svelte"
    import { Toaster } from 'svelte-french-toast';
    import Footer from "../components/footer.svelte";
    import Konami from 'konami';

		Konami(() => {
        document.getElementById("logo").src = "/til.gif"
    })
</script>

<Header />
<Toaster />

<main>
{@render children()}
</main>

<Footer />
