<script setup lang="ts">
import TestCase from './views/TestCase.vue'
import { ref } from 'vue'
import { ElTree } from 'element-plus'

interface Tree {
  id: string
  label: string
  children?: Tree[]
}

const testCaseName = ref('')
const testSuite = ref('')
const handleNodeClick = (data: Tree) => {
  testCaseName.value = data.label
  testSuite.value = data.parent
}

const data = ref([])
const treeRef = ref<InstanceType<typeof ElTree>>()

const requestOptions = {
    method: 'POST'
};
fetch('/server.Runner/GetSuites', requestOptions)
    .then(response => response.json())
    .then(d => {
      data.value = []
      Object.keys(d.data).map(k => {
        console.log(d.data[k])
        let suite = {
          id: k,
          label: k,
          children: [],
        }

        d.data[k].data.forEach((item: any) => {
          suite.children?.push({
            id: item,
            label: item,
            parent: k,
          })
        })
        data.value.push(suite)
      })

      // treeRef.value.updateKeyChildren('1', data[0].children)
    });

function load(n) {
  console.log(n)
}

const renderContent = (
  h,
  {
    node,
    data,
    store,
  }: {
    node: Node
    data: Tree
    store: Node['store']
  }
) => {
  if (node.childNodes.length > 0) {
    return h(
      'span',
      {
        class: 'custom-tree-node',
      },
      h('span', null, node.label),
      h(
        'span',
        null,
        h(
          'a',
          {
          },
          'Append'
        )
      )
    )
  } else {
    return h(
      'span',
      {
        class: 'custom-tree-node',
      },
      h('span', null, node.label)
    )
  }
}
</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-aside width="200px">
        <el-tree :data="data" :props="defaultProps"
          default-expand-all
          ref="treeRef"
          node-key="id"
          @node-click="handleNodeClick" />
      </el-aside>

      <el-main>
        <TestCase :suite="testSuite" :name="testCaseName"/>
      </el-main>
    </el-container>
  </div>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  margin-top: 2rem;
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;

    padding: 1rem 0;
    margin-top: 1rem;
  }
}
.demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}
</style>
