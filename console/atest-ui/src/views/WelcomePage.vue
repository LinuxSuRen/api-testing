<template>
  <!-- todo: add logo as bgc img -->
  <div class="index">
    <br />
    <el-card shadow="hover">
      <template #header>
        <h2>Welcome to use atest to improve your code quality!</h2>
      </template>
      <div>Please read the following guide if this is your first time to use atest.</div>
      <br />
      <li>Create a store for saving the data;</li>
      <li>Create a test suite on the left panel;</li>
      <li>Select a suite, then create the test case;</li>
      <br />
      <div>
        Please get more details from the
        <a href="https://linuxsuren.github.io/api-testing/" target="_blank" rel="noopener"
          >official document</a
        >.
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.index {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 80%;
}

li {
  padding-left: 20px;
}

h2 {
  font-style: italic;
  text-align: center;
}
</style>
