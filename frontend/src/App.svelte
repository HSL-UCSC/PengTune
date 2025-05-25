<script>
  import logo from './assets/images/HSLlogo.png'
  import 'round-slider/dist/roundslider.min.css'; // ✅ move here

  let resultText = "Please enter your name below 👇"
  let name

  function greet() {
    Greet(name).then(result => resultText = result)
  }

  import KnobPanel from './components/KnobPanel.svelte';
</script>

<main>
  <img alt="HSL logo" id="logo" src="{logo}">
  <h1>PID Tuning Panel</h1>
  <KnobPanel />
</main>

<style>

  #logo {
    display: block;
    width: 20%;
    height: 20%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }

</style>
