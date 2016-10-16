$module('funcref', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (value) {
      var instance = new $static();
      var init = [];
      instance.value = value;
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.SomeFunction = function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($this.value);
        return;
      };
      return $promise.new($continue);
    };
    this.$typesig = function () {
      return $t.createtypesig(['SomeFunction', 2, $g.____testlib.basictypes.Function($g.____testlib.basictypes.Boolean).$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.funcref.SomeClass).$typeref()]);
    };
  });

  $static.AnotherFunction = function (toCall) {
    var $result;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            toCall().then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
  $static.TEST = function () {
    var $result;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.funcref.SomeClass.new($t.box(true, $g.____testlib.basictypes.Boolean)).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            sc = $result;
            $g.funcref.AnotherFunction($t.dynamicaccess(sc, 'SomeFunction')).then(function ($result0) {
              $result = $result0;
              $current = 2;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 2:
            $resolve($result);
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  };
});
