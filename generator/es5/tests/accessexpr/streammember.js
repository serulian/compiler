$module('streammember', function () {
  var $static = this;
  this.$class('SomeClass', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      var init = [];
      init.push($promise.resolve($t.box(2, $g.____testlib.basictypes.Integer)).then(function (result) {
        instance.SomeInt = result;
      }));
      return $promise.all(init).then(function () {
        return instance;
      });
    };
  });

  $static.AnotherThing = function (somestream) {
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      $t.streamaccess(somestream, 'SomeInt');
      $resolve();
      return;
    };
    return $promise.new($continue);
  };
});
