$module('structuralcast', function () {
  var $static = this;
  this.$class('BaseClass', true, '', function (T) {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Result = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(true, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['Result', 3, $g.____testlib.basictypes.Boolean.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.structuralcast.BaseClass(T)).$typeref()]);
    };
  });

  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($g.structuralcast.BaseClass($g.____testlib.basictypes.Integer).new().then(function (value) {
        instance.BaseClass$Integer = value;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
    $instance.Result = $t.property(function () {
      var $this = this;
      var $current = 0;
      var $continue = function ($resolve, $reject) {
        $resolve($t.box(false, $g.____testlib.basictypes.Boolean));
        return;
      };
      return $promise.new($continue);
    });
    this.$typesig = function () {
      return $t.createtypesig(['Result', 3, $g.____testlib.basictypes.Boolean.$typeref()], ['new', 1, $g.____testlib.basictypes.Function($g.structuralcast.SomeClass).$typeref()]);
    };
  });

  $static.DoSomething = function (sc) {
    return $promise.empty();
  };
  $static.TEST = function () {
    var $result;
    var sc;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      while (true) {
        switch ($current) {
          case 0:
            $g.structuralcast.SomeClass.new().then(function ($result0) {
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
            sc.BaseClass$Integer.Result().then(function ($result0) {
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
